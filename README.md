# sway-kb-comfi
Simple util written in go lang for the Sway env (uses Sway IPC). It makes sure your Keyboard layout stays dedicated per each window.

## Problem I tried to solve
Working in other enviroments and OSs I made use of a great feature "dedicated keyboard layout per application". It allows me to have the Telegram opened and running with my native keyboard layout, while other programs (including those being launched as new) to have a default en-US layout until changed explicitly. Memorizing layout per running application makes switching between applications comfortable and does not require to switch back and forth constantly like when you have a single global keyboard layout.

## Solution
I've written a small tool that starts withing Sways's ```config``` file as simple as ```exec sway-kb-comfi``` provided the ```sway-kb-comfi``` tool is on your PATH

## How it works
The daemon listens to Sway IPC window focus events. Whenever the focus moves to another window:
1. The layout that was active up to that moment is persisted for the window being left.
1. The layout remembered for the window being focused is restored (or the default layout, if the window is seen for the first time). The switch only happens if the remembered layout actually differs from the active one.

Layouts are keyed by the Sway container ID, which is unique for the lifetime of a window (unlike the process PID, which the OS can reuse). State is kept in memory only.

Note: the tool intentionally does not react to `xkb_layout` input events by switching anything. With Sway's keyboard grouping those events are unreliable — grouped keyboards don't emit them, while auxiliary devices (headsets, mice with keyboard interfaces, ACPI hotkeys) emit spurious ones. Reacting to them used to make the layout flip back and forth by itself. They are only logged (use `-v` and `-keyboard-id` for diagnostics).

## Compilation
as simple as 

```bash
go mod tidy
go build
```

## Flags
* `-v` — verbose logging (without it, logs are discarded)
* `-keyboard-id <identifier>` — preferred keyboard identifier for log output only; layout events are accepted from all keyboard devices regardless, since Sway keyboard groups share one layout. Can be obtained from:
  `swaymsg -t get_inputs | jq -r '.[] | select(.type=="keyboard") | .identifier'`

## Running as a systemd user service
A user unit is provided at `systemd/sway-kb-comfi.service`. It expects the binary at `~/.local/bin/sway-kb-comfi` (`go build -o ~/.local/bin/sway-kb-comfi .`), starts with the graphical session and stops together with it. To use it:

```bash
ln -s "$PWD/systemd/sway-kb-comfi.service" ~/.config/systemd/user/
systemctl --user daemon-reload
systemctl --user enable --now sway-kb-comfi.service
```

Logs: `journalctl --user -u sway-kb-comfi.service` (rebuild with `-ldflags` or run manually with `-v` for verbose output).

## References
* This tool heavily relies on a wonderful package [go-sway](https://github.com/joshuarubin/go-sway) that introduces IPC bindings for Sway with the concept of custom user handlers. Amazing stuff and a good place for me to learn more on best-practices in go. Thanks to the author [joshuarubin](https://github.com/joshuarubin)!
* Analogue that also inspired me to write my own stuff is a [swaykbdd](https://github.com/artemsen/swaykbdd) written by [artemsen](https://github.com/artemsen) in C. However it supports even different layouts within one Application that uses tabbed interface, which is really nice to have if you are working in a Web Browser and have some web applications in tabs that need a dedicated layout. Some nice stuff, thanks for the inspiration!
