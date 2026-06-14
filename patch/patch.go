// Package patch is grove's batched DOM patch protocol. Instead of making one
// wasm↔JS call per DOM mutation, a Recorder collects mutations into a Buffer
// as the reconciler runs, Encode serializes the whole buffer to a single
// string, and Apply replays that string onto a target — in the browser, one
// JS-side interpreter applies the batch, turning N boundary crossings per
// commit into one.
//
// Nodes are referenced by small integer ids rather than live handles, so the
// op stream is self-contained and platform-neutral: the same encoding drives
// the host reference Apply (used in tests, replaying onto any
// renderer.Renderer) and the browser applier. The container is always id 1;
// created nodes are numbered from 2; id 0 means "none" (append, or no
// sibling).
package patch

import (
	"strconv"
	"strings"

	"github.com/gyoumi/grove/renderer"
)

// ContainerID is the reserved id of the mount container; pass it to
// grove.Mount as the container handle when rendering through a Recorder.
const ContainerID = 1

// opcodes — one byte each, the first field of every op.
const (
	opCreateElement = "e" // e id tag groveID
	opCreateText    = "t" // t id text
	opSetText       = "s" // s id text
	opSetAttr       = "a" // a id name value
	opRemoveAttr    = "r" // r id name
	opSetProp       = "p" // p id name kind value   (kind: s|b|n|z=nil)
	opInsert        = "i" // i parent child before
	opRemove        = "x" // x parent child
	opListen        = "l" // l event
)

// Buffer accumulates encoded ops. The wire format is one op per line,
// fields separated by tabs, with tabs and newlines in text escaped — compact,
// allocation-light, and trivial to parse on either side without a JSON
// dependency in the hot path.
type Buffer struct {
	b     strings.Builder
	count int
}

func (buf *Buffer) field(s string) {
	// Escape the field separators so arbitrary text/attribute values are safe.
	if strings.IndexByte(s, '\t') >= 0 || strings.IndexByte(s, '\n') >= 0 || strings.IndexByte(s, '\\') >= 0 {
		s = strings.NewReplacer("\\", "\\\\", "\t", "\\t", "\n", "\\n").Replace(s)
	}
	buf.b.WriteString(s)
}

func (buf *Buffer) op(code string, fields ...string) {
	if buf.count > 0 {
		buf.b.WriteByte('\n')
	}
	buf.count++
	buf.b.WriteString(code)
	for _, f := range fields {
		buf.b.WriteByte('\t')
		buf.field(f)
	}
}

// CreateElement records creating an element with the given renderer id, tag,
// and grove id (for event delegation).
func (buf *Buffer) CreateElement(id int, tag string, groveID int) {
	buf.op(opCreateElement, itoa(id), tag, itoa(groveID))
}

func (buf *Buffer) CreateText(id int, text string) { buf.op(opCreateText, itoa(id), text) }
func (buf *Buffer) SetText(id int, text string)    { buf.op(opSetText, itoa(id), text) }
func (buf *Buffer) SetAttr(id int, name, value string) {
	buf.op(opSetAttr, itoa(id), name, value)
}
func (buf *Buffer) RemoveAttr(id int, name string) { buf.op(opRemoveAttr, itoa(id), name) }

// SetProp records a property write. The value's kind is tagged so Apply can
// restore string/bool/number/nil rather than guessing from text.
func (buf *Buffer) SetProp(id int, name string, value any) {
	switch v := value.(type) {
	case nil:
		buf.op(opSetProp, itoa(id), name, "z", "")
	case string:
		buf.op(opSetProp, itoa(id), name, "s", v)
	case bool:
		buf.op(opSetProp, itoa(id), name, "b", strconv.FormatBool(v))
	case int:
		buf.op(opSetProp, itoa(id), name, "n", strconv.Itoa(v))
	case float64:
		buf.op(opSetProp, itoa(id), name, "n", strconv.FormatFloat(v, 'g', -1, 64))
	default:
		// Fall back to the value's default string form, tagged as a string.
		buf.op(opSetProp, itoa(id), name, "s", toString(v))
	}
}

func (buf *Buffer) Insert(parent, child, before int) {
	buf.op(opInsert, itoa(parent), itoa(child), itoa(before))
}
func (buf *Buffer) Remove(parent, child int) { buf.op(opRemove, itoa(parent), itoa(child)) }
func (buf *Buffer) Listen(event string)      { buf.op(opListen, event) }

// Len is the number of ops recorded since the last Reset.
func (buf *Buffer) Len() int { return buf.count }

// Encode returns the recorded ops as a single string.
func (buf *Buffer) Encode() string { return buf.b.String() }

// Reset clears the buffer for the next batch.
func (buf *Buffer) Reset() {
	buf.b.Reset()
	buf.count = 0
}

func itoa(i int) string { return strconv.Itoa(i) }

func toString(v any) string {
	type stringer interface{ String() string }
	if s, ok := v.(stringer); ok {
		return s.String()
	}
	return ""
}

// Recorder is a renderer.Renderer that records mutations into a Buffer
// instead of touching a real DOM; node handles are ints. It is the shared
// recording half of a batched renderer — the browser renderer embeds one and
// supplies a real Flush and event delegation, while tests drive it directly
// and replay its flushes with Apply. Event reads and Listen are no-ops here;
// platform renderers override them.
type Recorder struct {
	Buf      Buffer
	nextID   int
	dispatch renderer.Dispatch

	pending []func()
	flushes []string // encoded batches captured by Flush (for tests)
}

// NewRecorder returns a Recorder whose container is ContainerID.
func NewRecorder() *Recorder { return &Recorder{nextID: ContainerID + 1} }

func (r *Recorder) id(n renderer.Node) int {
	if n == nil {
		return 0
	}
	return n.(int)
}

func (r *Recorder) SetDispatch(d renderer.Dispatch) { r.dispatch = d }

// PortalRoot is the container id; the batched browser renderer inherits this,
// so portal children attach under the mount container.
func (r *Recorder) PortalRoot() renderer.Node { return ContainerID }

func (r *Recorder) CreateElement(tag string, groveID int) renderer.Node {
	id := r.nextID
	r.nextID++
	r.Buf.CreateElement(id, tag, groveID)
	return id
}

func (r *Recorder) CreateText(text string) renderer.Node {
	id := r.nextID
	r.nextID++
	r.Buf.CreateText(id, text)
	return id
}

func (r *Recorder) SetText(n renderer.Node, text string)        { r.Buf.SetText(r.id(n), text) }
func (r *Recorder) SetAttr(n renderer.Node, name, value string) { r.Buf.SetAttr(r.id(n), name, value) }
func (r *Recorder) RemoveAttr(n renderer.Node, name string)     { r.Buf.RemoveAttr(r.id(n), name) }
func (r *Recorder) SetProp(n renderer.Node, name string, value any) {
	r.Buf.SetProp(r.id(n), name, value)
}
func (r *Recorder) InsertBefore(parent, child, before renderer.Node) {
	r.Buf.Insert(r.id(parent), r.id(child), r.id(before))
}
func (r *Recorder) Remove(parent, child renderer.Node) { r.Buf.Remove(r.id(parent), r.id(child)) }
func (r *Recorder) Listen(event string)                { r.Buf.Listen(event) }

// Flush captures the current batch (for tests to replay) and resets the
// buffer. Platform renderers override this to ship the batch across the
// boundary.
func (r *Recorder) Flush() {
	if r.Buf.Len() == 0 {
		return
	}
	r.flushes = append(r.flushes, r.Buf.Encode())
	r.Buf.Reset()
}

func (r *Recorder) Schedule(f func()) { r.pending = append(r.pending, f) }

// Settle runs scheduled work until none remains; tests use it to drive
// update flushes the way testdom.R.Settle does.
func (r *Recorder) Settle() {
	for i := 0; len(r.pending) > 0; i++ {
		if i > 1000 {
			panic("patch: scheduled work never settled")
		}
		q := r.pending
		r.pending = nil
		for _, f := range q {
			f()
		}
	}
}

// Flushes returns the encoded batches captured so far (one per commit).
func (r *Recorder) Flushes() []string { return r.flushes }

// Dispatch delivers an event to the reconciler (test helper).
func (r *Recorder) Dispatch(id int, event string, raw any) {
	if r.dispatch != nil {
		r.dispatch(id, event, raw)
	}
}

// Event reads are unused without a platform event object.
func (r *Recorder) PreventDefault(any)         {}
func (r *Recorder) StopPropagation(any)        {}
func (r *Recorder) Str(any, ...string) string  { return "" }
func (r *Recorder) Bool(any, ...string) bool   { return false }
func (r *Recorder) Num(any, ...string) float64 { return 0 }
