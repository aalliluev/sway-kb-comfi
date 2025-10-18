package helpers

import (
	"aalliluev/sway-kb-comfi/handlers"
	"context"
	"log"
	"strings"

	"github.com/joshuarubin/go-sway"
)

func CheckKeyboardAndFindBestIfNotSpecified(ctx context.Context, client sway.Client, keyboardId string) string {

	devices := getKeyboards(ctx, client)

	if keyboardId != "" {
		// check if it does exist indeed in a current SWAY setup
		log.Printf("Trying to find device specifed as %s", keyboardId)
		if exists := getKeyboardExactMatch(keyboardId, devices); exists {
			log.Println("Found")
			return keyboardId
		}
		log.Printf("The device specified as '%s' was not found. I'll scan and pick an appropriate one automatically", keyboardId)
	} else {
		log.Printf("Making best guess to determine a keyboard...")
	}

	// need to determine a single Input device we will be listening changes from.
	result := getKeyboardBestMatch(devices)

	if result == "" {
		log.Fatalln("Could not find any suitable input device")
	}
	return result
}

func getKeyboards(ctx context.Context, client sway.Client) *map[string]string {
	result := make(map[string]string)

	inputs, err := client.GetInputs(ctx)
	if err != nil {
		log.Fatalln("Couldn't get input devices from SWAY IPC.", err)
	}
	for _, i := range inputs {
		isKeyboard := i.Type == handlers.INPUT_TYPE_KEYBOARD
		if isKeyboard {
			result[i.Identifier] = i.Identifier
		}
	}

	return &result
}

func getKeyboardExactMatch(keyboard string, devices *map[string]string) bool {
	for i := range *devices {
		if i == keyboard {
			return true
		}
	}
	return false
}

func getKeyboardBestMatch(devices *map[string]string) string {
	worstCaseScenarioChoice := ""
	for i := range *devices {
		if worstCaseScenarioChoice != "" {
			worstCaseScenarioChoice = i
		}
		containsWordKeyboard := strings.Contains(strings.ToLower(i), handlers.INPUT_TYPE_KEYBOARD)
		if containsWordKeyboard {
			return i
		}
	}
	return worstCaseScenarioChoice
}
