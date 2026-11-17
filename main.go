package main

import (
	"aalliluev/sway-kb-comfi/handlers"
	"aalliluev/sway-kb-comfi/helpers"
	"aalliluev/sway-kb-comfi/layouts"
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joshuarubin/go-sway"
)

func main() {

	// some flags
	verbose := flag.Bool("v", false, "adds verbosity")
	keyboardID := flag.String("keyboard-id", "", "Keyboard identifier to report in logs. Layout events are accepted from\nall keyboard devices regardless, since sway keyboard groups share the layout.\nIdentifier can be obtained from command:\nswaymsg -t get_inputs | jq -r '.[] | select(.type==\"keyboard\") | .identifier'")
	flag.Parse()

	if !*verbose {
		log.SetOutput(io.Discard)
	}

	socket := os.ExpandEnv("${SWAYSOCK}")
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	defer cancel()

	// client for commands
	client, err := sway.New(ctx, sway.WithSocketPath(socket))
	if err != nil {
		log.Fatalln("Failed to create IPC client", err)
	}

	if *keyboardID != "" {
		*keyboardID = helpers.CheckKeyboardAndFindBestIfNotSpecified(ctx, client, *keyboardID)
	} else {
		log.Printf("No keyboard specified with -keyboard-id; layout events from all keyboards will be accepted")
	}

	log.Printf("Preferred input device is: %s", *keyboardID)

	// main logic for handling Sway Events in a way to persist and restore layouts
	handler := handlers.CreatelayoutHandler(
		client,
		&layouts.GlobalLayouts{
			DefaultLayout: 0,
			Layouts:       make(map[int64]*layouts.LayoutInfo),
		},
		*keyboardID,
	)

	// memorize the currently focused window so a restart doesn't cause a layout switch
	handler.Seed(ctx)

	// subscribing to events
	go func() {
		for {
			err := sway.Subscribe(ctx, handler, sway.EventTypeWindow, sway.EventTypeInput)
			if err == nil {
				return
			}
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("Failed to subscribe: %v. Retrying in 5s...", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
		}
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// all good we're back in business
	log.Println("--- WORKING ---")

	sig := <-sigs

	log.Printf("%v received. Terminating...", sig)

	cancel()
	os.Exit(0)
}
