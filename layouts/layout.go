package layouts

import (
	"log"
)

const (
	INVALID_LAYOUT int = -1 // to identify unset layout
)

type LayoutInfo struct {
	ID     int64  // sway container ID (unique for the lifetime of the container)
	Layout int    // current layout
	Name   string // Either Name or Class what's available
}

type GlobalLayouts struct {
	Layouts       map[int64]*LayoutInfo // persisted state per sway container ID
	DefaultLayout int                   // default layout to be used for those windows that are not yet tracked
}

func (gl *GlobalLayouts) GetPersistedLayout(id int64) int {
	if l, ok := gl.Layouts[id]; ok {
		return l.Layout
	}
	return INVALID_LAYOUT
}

func (gl *GlobalLayouts) PersistLayoutForId(id int64, name string, layout int) {
	l, ok := gl.Layouts[id]
	if !ok {
		l = &LayoutInfo{ID: id, Name: name}
		gl.Layouts[id] = l
	}
	l.Layout = layout
	log.Printf("Persisted layout %d for container ID: %d (%s)", layout, id, name)
}
