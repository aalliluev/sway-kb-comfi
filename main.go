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

	"github.com/joshuarubin/go-sway"
)

func main() {

	// some flags
	verbose := flag.Bool("v", false, "adds verbosity")
	defaultInputIdentifier := flag.String("keyboard-id", "", "Default input identifier to capture events from. Can be obtained from command\nswaymsg -t get_inputs | jq -r '.[] | select(.type==\"keyboard\") | .identifier'\nIf not specified, I'll do my best to infer the first one automatically")
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

	*defaultInputIdentifier = helpers.CheckKeyboardAndFindBestIfNotSpecified(ctx, client, *defaultInputIdentifier)

	log.Printf("Input device to listen to is: %s", *defaultInputIdentifier)

	// main logic for handling Sway Events in a way to persist and restore layouts
	handler := handlers.CreatelayoutHandler(
		client,
		&layouts.GlobalLayouts{
			DefaultLayout: 0,
			Layouts:       make(map[uint32]*layouts.LayoutInfo),
		},
		*defaultInputIdentifier, //"5426:674:Razer_Razer_Ornata_V3_X_Keyboard",
	)

	// subscribing to events
	go func() {
		if err = sway.Subscribe(ctx, handler, sway.EventTypeWindow, sway.EventTypeInput); err != nil && !errors.Is(err, context.Canceled) {
			log.Fatalln("Failed to subscribe.", err)
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
