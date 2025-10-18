package layouts

import (
	"log"
)

const (
	INVALID_LAYOUT int = -1 // to identified unset layout
)

type LayoutInfo struct {
	Pid    uint32 // Process Id
	Layout int    // current layout
	Name   string // Either Name or Class what's available
}

type GlobalLayouts struct {
	Layouts       map[uint32]*LayoutInfo // persisted state per PID. TODO: dedicated state per other parameters? like for Tabbed Applications (browser)?
	DefaultLayout int                    // default layout to be used for those processes that are not yet tracked (new too)
}

func (gl *GlobalLayouts) GetOrCreate(pid uint32, name string) (*LayoutInfo, bool) {
	val, exists := gl.Layouts[pid]
	if !exists {
		val = &LayoutInfo{Pid: pid, Layout: gl.DefaultLayout, Name: name}
		gl.Layouts[pid] = val
		return val, false

	}
	return val, true
}

func (gl *GlobalLayouts) GetPersistedLayout(id uint32, name string) (int, bool) {
	l, existed := gl.GetOrCreate(id, name)
	if l != nil {
		// returning either existing Layout or Default one accoring to settings.
		if existed {
			log.Printf("Restored a persisted layout %d for ID: %d", l.Layout, id)
			return l.Layout, true
		} else {
			log.Printf("Restoring Default layout %d for: PID: %d", gl.DefaultLayout, id)
			return gl.DefaultLayout, false
		}
	}

	// if we are here returning INVALID_LAYOUT in order to prevent restoring an invalid value for this PID.
	return INVALID_LAYOUT, false
}

func (gl *GlobalLayouts) PersistLayoutForId(id uint32, newLayoutIndex int) bool {
	l, _ := gl.GetOrCreate(id, "")

	if l != nil {
		l.Layout = newLayoutIndex
		log.Printf("Persisted layout %d for: ID: %d", newLayoutIndex, id)
		return true
	}

	return false
}
