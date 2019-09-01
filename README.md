Build
=====

Just use `make`.

The binaries are put into the `build` directory.

| File                  | Platform              |
|-----------------------|-----------------------|
| lightbull-x64-linux   | Obviously             |
| lightbull-armv7-linux | Raspberry Pi 3        |
| lightbull-armv5-linux | Older Raspberry Pis   |

Customize LEDs
==============

Concept
-------

The software needs to know, which LED IDs belong to which parts. This is defined in `hardware/hardware.go`:

    hw.Led.AddPart("horn_left", 20, 10)
    hw.Led.AddPart("horn_left", 40, 44)
    hw.Led.AddPart("horn_left", 70, 75)

    hw.Led.AddPart("horn_right", 0, 5)

The LED IDs from the first to second number (including both IDs!) is added to that part.
So here, the part "horn_left" would consist of these IDs:

    20, 19, 18, ... 11, 10, 40, 41, ... 44, 70, 71, .... 75

Calibrate
---------

`lightbull-arch-os calibrate` allows to switch interactively single LEDs on and may be helpful to find out,
which LEDs belong to which part.

By default, the tool sends out control commands for 750 LEDs. If this number if to low, it can be adjusted with
the `-n` parameter.

Test
----

`lightbull-arch-os test` runs a small test program.

Control server
==============

`lightbull-arch-os run` runs the control server, the web interface is accessible on port 8080.

Some settings can be changed using the configuration file which can be places in `/etc/lightbull/config.yaml` or `./config.yaml`.
