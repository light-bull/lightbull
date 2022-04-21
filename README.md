# lightbull

![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/light-bull/lightbull/Build/main?style=plastic)
![GitHub last commit (branch)](https://img.shields.io/github/last-commit/light-bull/lightbull/main?style=plastic)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/light-bull/lightbull?style=plastic)

## Build

Just use `make`.

The binaries are put into the `build` directory.

| File                  | Platform              |
|-----------------------|-----------------------|
| lightbull-x64-linux   | Obviously             |
| lightbull-armv7-linux | Raspberry Pi 3        |
| lightbull-armv5-linux | Older Raspberry Pis   |

## Customize LEDs

### Concept

The software needs to know, which LED IDs belong to which parts. This is defined in the config file `config.yaml`:

    leds:
        parts:
            - name: "horn_right"
              leds: [[0, 5]]
            - name: "horn_left"
              leds: [[20, 10], [40, 44], [70, 75]]

The LED IDs from the first to second number (including both IDs!) is added to that part.
So here, the part "horn_left" would consist of these IDs:

    20, 19, 18, ... 11, 10, 40, 41, ... 44, 70, 71, .... 75

### Calibrate

`lightbull-arch-os calibrate` allows to switch interactively single LEDs on and may be helpful to find out,
which LEDs belong to which part.

By default, the tool sends out control commands for 750 LEDs. If this number if to low, it can be adjusted with
the `-n` parameter.

### Test

`lightbull-arch-os test` runs a small test program.

## Control server

`lightbull-arch-os run` runs the control server, the API is accessible on port 8080.

Some settings can be changed using the configuration file which can be places in `/etc/lightbull/config.yaml` or `./config.yaml`.

## Development

### Code checks

We use pre-commit for code and styleguide checks.

Install it once as git hook:

    pre-commit install

Run pre-commit manually:

    pre-commit run --all-files

### Add new data type for parameters

* Add name to `shows/parameters/const.go`
* Create new datatype in `shows/parameters/....go` based on existing one
* Add in `NewParameter` function in `shows/parameters/parameter.go`

### Add a new effect

* Add name to `shows/effects/const.go`
* Add to `GetEffects` in same file
* Create new effect in `shows/effects/....go` based on existing one
* Add in `NewEffect` function in `shows/effects/effect.go`
