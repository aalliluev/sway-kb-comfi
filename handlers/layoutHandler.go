package handlers

import (
	"aalliluev/sway-kb-comfi/layouts"
	"context"
	"fmt"
	"log"

	"github.com/joshuarubin/go-sway"
)

const (
	INVALID_PID             uint32 = 0xFFFFFFFF
	INPUT_TYPE_KEYBOARD     string = "keyboard"
	INPUT_EVENT_CHANGE_NAME string = "xkb_layout"
)

type PersistableLayoutHandler struct {
	sway.EventHandler
	client                 sway.Client
	TrackingLayouts        *layouts.GlobalLayouts
	DefaultInputIdentifier string
	LastFocusedId          uint32
}

func CreatelayoutHandler(client sway.Client, layouts *layouts.GlobalLayouts, defaultKeyboardName string) *PersistableLayoutHandler {
	return &PersistableLayoutHandler{
		client:                 client,
		TrackingLayouts:        layouts,
		LastFocusedId:          INVALID_PID,
		DefaultInputIdentifier: defaultKeyboardName,
	}
}

// SWAY IPC Handler for WindowEvent type
func (h *PersistableLayoutHandler) Window(ctx context.Context, e sway.WindowEvent) {
	switch e.Change {

	case "focus":
		focused := e.Container.FocusedNode()
		if focused == nil {
			return
		}
		name := focused.Name
		id := *focused.PID
		class := ""
		if focused.WindowProperties != nil {
			class = focused.WindowProperties.Class
		}
		log.Printf("Got focused: PID: %d, Name: %s, Class: %s", id, name, class)
		if h.LastFocusedId != id {
			h.LastFocusedId = id
			if l, _ := h.TrackingLayouts.GetPersistedLayout(id, name); l != layouts.INVALID_LAYOUT {
				h.changeLayout(ctx, l)
			}
		}
	case "close":
		delete(h.TrackingLayouts.Layouts, *e.Container.PID)
	}
}

// SWAY IPC Handler for InputEvent type
func (h *PersistableLayoutHandler) Input(ctx context.Context, e sway.InputEvent) {
	filteredId := e.Input.Identifier == h.DefaultInputIdentifier
	filteredByType := e.Input.Type == INPUT_TYPE_KEYBOARD
	filteredByChange := e.Change == INPUT_EVENT_CHANGE_NAME

	if !filteredByType || !filteredId || !filteredByChange {
		return
	}

	input := e.Input
	id := input.Identifier
	l := *input.XKBActiveLayoutIndex
	tp := input.Type
	log.Printf("Layoutchanged, persisting for ID: %s, Type: %s, LayoutIndex: %d", id, tp, l)
	if ok := h.TrackingLayouts.PersistLayoutForId(h.LastFocusedId, int(l)); ok {
		h.changeLayout(ctx, int(l))
	}
}

// helper method to send command to SWAY
func (h *PersistableLayoutHandler) changeLayout(ctx context.Context, layoutIndex int) bool {
	replies, err := h.client.RunCommand(ctx, fmt.Sprintf("input * xkb_switch_layout %d", layoutIndex))
	if err != nil {
		log.Printf("Failed to run command\n%v", err)
		return false
	}

	if !replies[0].Success {
		log.Printf("Failed to execute command. Sway Error: %s", replies[0].Error)
		return false
	}
	return true
}
