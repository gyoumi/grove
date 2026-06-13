package patch_test

import (
	"strconv"
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/patch"
	"github.com/gyoumi/grove/testdom"
)

var probeSetter func([]string)

// protoApp exercises the protocol surface: elements, text, classes and
// attributes, a keyed list (insert/remove/reorder → moves), a value property
// (SetProp), and text that changes on update.
func protoApp() *g.Node {
	items, set := g.UseState([]string{"a", "b", "c"})
	probeSetter = set
	return g.Div(g.Class("wrap"),
		g.Ul(g.Map(items, func(s string) *g.Node {
			return g.Li(g.Key(s), g.Class("item"), g.Attr("data-k", s), s)
		})),
		g.Input(g.Value(strconv.Itoa(len(items)))),
		g.P("count ", strconv.Itoa(len(items))),
	)
}

// A render through the batched protocol — recorded as an op buffer, encoded,
// and replayed by Apply onto a fresh renderer — must reconstruct exactly the
// tree a direct render produces, at mount and after every kind of update.
func TestPatchProtocolMatchesDirectRender(t *testing.T) {
	direct := testdom.Mount(g.C0(protoApp))
	setDirect := probeSetter

	rec := patch.NewRecorder()
	g.Mount(rec, patch.ContainerID, g.C0(protoApp))
	setRec := probeSetter

	// Replay the recorded batches onto a fresh renderer.
	applied := testdom.New()
	ap := patch.NewApplier(applied, applied.Container)
	for _, batch := range rec.Flushes() {
		ap.Apply(batch)
	}
	if applied.HTML() != direct.HTML() {
		t.Fatalf("mount mismatch:\n direct:  %s\n applied: %s", direct.HTML(), applied.HTML())
	}

	updates := [][]string{
		{"a", "b", "c", "d"},      // append
		{"d", "a", "b", "c"},      // rotate (one move)
		{"a", "c"},                // remove two
		{"x", "a", "y", "c", "z"}, // insert three, reorder
		{"z", "y", "x"},           // shrink + reverse
		{"a", "b", "c"},           // back to start
	}
	for _, u := range updates {
		setDirect(u)
		direct.Settle()

		seen := len(rec.Flushes())
		setRec(u)
		rec.Settle()
		for _, batch := range rec.Flushes()[seen:] {
			ap.Apply(batch)
		}

		if applied.HTML() != direct.HTML() {
			t.Fatalf("update %v mismatch:\n direct:  %s\n applied: %s", u, direct.HTML(), applied.HTML())
		}
	}
}

// Buffer escaping must survive text containing the field and op separators
// (tabs and newlines) and backslashes.
func TestBufferEscapesSeparators(t *testing.T) {
	tricky := "a\tb\nc\\d"
	applied := testdom.New()

	var buf patch.Buffer
	buf.CreateElement(2, "p", 7)
	buf.CreateText(3, tricky)
	buf.Insert(2, 3, 0) // text into the <p>
	buf.Insert(patch.ContainerID, 2, 0)
	patch.Apply(buf.Encode(), applied, applied.Container)

	if got := applied.Container.Children[0].Children[0].Text; got != tricky {
		t.Fatalf("escaped text round-trip: got %q want %q", got, tricky)
	}
}
