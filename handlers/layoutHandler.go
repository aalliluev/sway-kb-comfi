package handlers

import (
	"aalliluev/sway-kb-comfi/layouts"
	"context"
	"fmt"
	"log"

	"github.com/joshuarubin/go-sway"
)

const (
	INVALID_CONTAINER_ID    int64  = -1
	INPUT_TYPE_KEYBOARD     string = "keyboard"
	INPUT_EVENT_CHANGE_NAME string = "xkb_layout"
)

type PersistableLayoutHandler struct {
	sway.EventHandler
	client                 sway.Client
	TrackingLayouts        *layouts.GlobalLayouts
	DefaultInputIdentifier string
	LastFocusedConID       int64
	LastFocusedName        string
}

func CreatelayoutHandler(client sway.Client, layouts *layouts.GlobalLayouts, defaultKeyboardName string) *PersistableLayoutHandler {
	return &PersistableLayoutHandler{
		client:                 client,
		TrackingLayouts:        layouts,
		LastFocusedConID:       INVALID_CONTAINER_ID,
		DefaultInputIdentifier: defaultKeyboardName,
	}
}

// Seed memorizes the currently focused window so that a freshly started
// daemon doesn't treat the first focus event as a window change.
func (h *PersistableLayoutHandler) Seed(ctx context.Context) {
	tree, err := h.client.GetTree(ctx)
	if err != nil {
		log.Printf("Couldn't get tree: %v", err)
		return
	}
	if tree == nil {
		return
	}
	focused := tree.FocusedNode()
	if focused == nil || focused.PID == nil {
		return
	}
	h.LastFocusedConID = focused.ID
	h.LastFocusedName = focused.Name
	log.Printf("Seeded focused window: ConID: %d, Name: %s", focused.ID, focused.Name)
}

// SWAY IPC Handler for WindowEvent type
func (h *PersistableLayoutHandler) Window(ctx context.Context, e sway.WindowEvent) {
	switch e.Change {
	case "focus":
		h.handleFocus(ctx, e)
	case "close":
		delete(h.TrackingLayouts.Layouts, e.Container.ID)
		if h.LastFocusedConID == e.Container.ID {
			h.LastFocusedConID = INVALID_CONTAINER_ID
			h.LastFocusedName = ""
		}
	}
}

func (h *PersistableLayoutHandler) handleFocus(ctx context.Context, e sway.WindowEvent) {
	focused := e.Container.FocusedNode()
	if focused == nil || focused.PID == nil {
		// Not an application window (e.g. an empty workspace got focused).
		return
	}
	id := focused.ID
	if h.LastFocusedConID == id {
		return
	}
	name := focused.Name
	class := ""
	if focused.WindowProperties != nil {
		class = focused.WindowProperties.Class
	}
	log.Printf("Got focused: ConID: %d, Name: %s, Class: %s", id, name, class)

	// The layout is group-global in sway, so whatever is active right now is
	// the layout the previously focused window was used with.
	current, err := h.getCurrentLayout(ctx)
	if err != nil {
		log.Printf("Couldn't determine current layout, skipping: %v", err)
		return
	}

	if h.LastFocusedConID != INVALID_CONTAINER_ID {
		h.TrackingLayouts.PersistLayoutForId(h.LastFocusedConID, h.LastFocusedName, current)
	}
	h.LastFocusedConID = id
	h.LastFocusedName = name

	target := h.TrackingLayouts.GetPersistedLayout(id)
	if target == layouts.INVALID_LAYOUT {
		target = h.TrackingLayouts.DefaultLayout
	}
	h.changeLayout(ctx, target, current)
}

// SWAY IPC Handler for InputEvent type.
// Events are informational only: sway keyboard groups make these events
// unreliable (grouped keyboards don't emit them, pseudo-devices do), so the
// handler never sends commands from here — that used to cause a
// switch-layout-event feedback loop.
func (h *PersistableLayoutHandler) Input(ctx context.Context, e sway.InputEvent) {
	if e.Input.Type != INPUT_TYPE_KEYBOARD || e.Change != INPUT_EVENT_CHANGE_NAME || e.Input.XKBActiveLayoutIndex == nil {
		return
	}
	log.Printf("Layout changed on %s, index: %d (tracked keyboard: %s)",
		e.Input.Identifier, *e.Input.XKBActiveLayoutIndex, h.DefaultInputIdentifier)
}

func (h *PersistableLayoutHandler) getCurrentLayout(ctx context.Context) (int, error) {
	inputs, err := h.client.GetInputs(ctx)
	if err != nil {
		return layouts.INVALID_LAYOUT, err
	}
	for _, i := range inputs {
		if i.Type == INPUT_TYPE_KEYBOARD && i.XKBActiveLayoutIndex != nil {
			return int(*i.XKBActiveLayoutIndex), nil
		}
	}
	return layouts.INVALID_LAYOUT, fmt.Errorf("no keyboard input device found")
}

// helper method to send command to SWAY
func (h *PersistableLayoutHandler) changeLayout(ctx context.Context, layoutIndex, current int) {
	if layoutIndex == current {
		log.Printf("Layout %d already active, nothing to do", layoutIndex)
		return
	}

	replies, err := h.client.RunCommand(ctx, fmt.Sprintf("input * xkb_switch_layout %d", layoutIndex))
	if err != nil {
		log.Printf("Failed to run command\n%v", err)
		return
	}

	if len(replies) == 0 {
		log.Printf("Failed to run command: no reply from sway")
		return
	}

	if !replies[0].Success {
		log.Printf("Failed to execute command. Sway Error: %s", replies[0].Error)
		return
	}
	log.Printf("Switched layout to %d", layoutIndex)
}
