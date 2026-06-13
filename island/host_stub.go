//go:build !(js && wasm)

package island

import g "github.com/gyoumi/grove"

// Host receives island lifecycle calls outside the browser, where there is
// no JS runtime: el is the renderer handle for the container (a testdom
// node under go test). The default host is a no-op; tests install a
// recording one with SetHost to observe the calls grove would make.
type Host interface {
	Mount(el any, name, propsJSON string)
	Update(el any, name, propsJSON string)
	Unmount(el any, name string)
}

var host Host

// SetHost routes island lifecycle calls to h; nil restores the no-op.
func SetHost(h Host) { host = h }

func hostMount(ref *g.DOMRef, name, propsJSON string) {
	if host != nil {
		host.Mount(ref.Current, name, propsJSON)
	}
}

func hostUpdate(ref *g.DOMRef, name, propsJSON string) {
	if host != nil {
		host.Update(ref.Current, name, propsJSON)
	}
}

func hostUnmount(ref *g.DOMRef, name string) {
	if host != nil {
		host.Unmount(ref.Current, name)
	}
}
