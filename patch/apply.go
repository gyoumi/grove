package patch

import (
	"strconv"
	"strings"

	"github.com/gyoumi/grove/renderer"
)

// Apply replays an encoded batch onto target, mapping ids to the nodes target
// produces. container is registered as ContainerID up front. The id→Node map
// persists across calls when the same nodes is supplied, so successive batches
// (mount, then updates) build on each other; pass the same map back in via
// Applier for incremental application.
//
// This is the reference implementation of the protocol — the browser applier
// mirrors it in JS. Tests replay a Recorder's flushes onto a fresh testdom
// renderer and check the result matches a direct render.
func Apply(encoded string, target renderer.Renderer, container renderer.Node) {
	NewApplier(target, container).Apply(encoded)
}

// Applier holds the id→Node map so a sequence of batches can be applied
// incrementally onto the same target.
type Applier struct {
	target renderer.Renderer
	nodes  map[int]renderer.Node
}

// NewApplier returns an Applier whose container is ContainerID.
func NewApplier(target renderer.Renderer, container renderer.Node) *Applier {
	return &Applier{target: target, nodes: map[int]renderer.Node{ContainerID: container}}
}

func (ap *Applier) node(id int) renderer.Node {
	if id == 0 {
		return nil
	}
	return ap.nodes[id]
}

// Apply replays one encoded batch.
func (ap *Applier) Apply(encoded string) {
	if encoded == "" {
		return
	}
	for _, line := range strings.Split(encoded, "\n") {
		f := splitFields(line)
		switch f[0] {
		case opCreateElement:
			ap.nodes[atoi(f[1])] = ap.target.CreateElement(f[2], atoi(f[3]))
		case opCreateText:
			ap.nodes[atoi(f[1])] = ap.target.CreateText(f[2])
		case opSetText:
			ap.target.SetText(ap.node(atoi(f[1])), f[2])
		case opSetAttr:
			ap.target.SetAttr(ap.node(atoi(f[1])), f[2], f[3])
		case opRemoveAttr:
			ap.target.RemoveAttr(ap.node(atoi(f[1])), f[2])
		case opSetProp:
			ap.target.SetProp(ap.node(atoi(f[1])), f[2], decodeValue(f[3], f[4]))
		case opInsert:
			ap.target.InsertBefore(ap.node(atoi(f[1])), ap.node(atoi(f[2])), ap.node(atoi(f[3])))
		case opRemove:
			ap.target.Remove(ap.node(atoi(f[1])), ap.node(atoi(f[2])))
		case opListen:
			ap.target.Listen(f[1])
		}
	}
}

// splitFields splits a tab-separated op line and unescapes each field.
func splitFields(line string) []string {
	raw := strings.Split(line, "\t")
	for i, s := range raw {
		if strings.IndexByte(s, '\\') >= 0 {
			s = strings.NewReplacer("\\t", "\t", "\\n", "\n", "\\\\", "\\").Replace(s)
		}
		raw[i] = s
	}
	return raw
}

func decodeValue(kind, v string) any {
	switch kind {
	case "z":
		return nil
	case "b":
		return v == "true"
	case "n":
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
		return 0
	default:
		return v
	}
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
