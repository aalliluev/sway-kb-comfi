# sway-kb-comfi
Simple util written in go lang for the Sway env (uses Sway IPC). It makes sure your Keyboard layout stays dedicated per each process.

## Problem I tried to solve
Working in other enviroments and OSs I made use of a great feature "dedicated keyboard layout per application". It allows me to have the Telegram opened and running with my native keyboard layout, while other programs (including those being launched as new) to have a default en-US layout until changed explicitly. Memorizing layout per running application makes switching between applications comfortable and does not require to switch back and forth constantly like when you have a single global keyboard layout.

## Solution
I've written a small tool that starts withing Sways's ```config``` file as simple as ```exec sway-kb-comfi``` provided the ```sway-kb-comfi``` tool is on your PATH

## Compilation
as simple as 

```bash
go mod tidy
go build
```

## References
* This tool heavily relies on a wonderful package [go-sway](https://github.com/joshuarubin/go-sway) that introduces IPC bindings for Sway with the concept of custom user handlers. Amazing stuff and a good place for me to learn more on best-practices in go. Thanks to the author [joshuarubin](https://github.com/joshuarubin)!
* Analogue that also inspired me to write my own stuff is a [swaykbdd](https://github.com/artemsen/swaykbdd) written by [artemsen](https://github.com/artemsen) in C. However it supports even different layouts within one Application that uses tabbed interface, which is really nice to have if you are working in a Web Browser and have some web applications in tabs that need a dedicated layout. Some nice stuff, thanks for the inspiration!
