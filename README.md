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
guild build
```
