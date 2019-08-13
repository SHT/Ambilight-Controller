# Ambilight

This project aims to deliver a robust client-server command-line application that is able to stream live, low-latency LED color data to a controller via a wireless / wired network connection.

## Dependencies

This project depends on the following libraries:
 - [cretz](https://github.com/cretz) / [go-scrap](https://github.com/cretz/go-scrap) - A wrapper around the Rust [scrap](https://github.com/quadrupleslap/scrap) library, which enables us to capture screenshot.
 - [jgarff](https://github.com/jgarff) / [rpi_ws281x](https://github.com/jgarff/rpi_ws281x) - Raspberry Pi library for controlling WS281X LEDs

## WIP

The project is currently under development. There might be bugs or inaccuracies on some parts of the documentation.

**TODO:**
- [ ] Build instructions
- [ ] Add dependencies on README
- [x] Setup instructions
- [x] Include more modes (Rainbow, Pulse etc.)
- [x] Documentation
- [x] Contributing instructions
- [x] Add binaries for Windows / Linux



## How it works

The *client* is launched on a Windows / Mac / Linux PC and the *server* on a controller (e.g. Raspberry Pi) that is wired with the LED strip that will be controlled.

Once the server is running on the controller, whenever the client is launched it will try to connect with the command-line parameters provided.

When the connection is established, the client will stream data to the server via a TCP socket, and the server will act according to the specified mode of operation (for now only Ambilight and Rainbow modes are supported, more will be added soon)

*Whatever dude how do I install it?*

## Installation instructions

#### Client (e.g. Windows PC)

1. Download the *client* binary for your operating system from [here](https://github.com/SHT/Ambilight/releases/latest/).
2. Create a `run.bat` file on the same folder as the binary.
3. Paste the following on the `run.bat` file, replacing `IP`, `PORT`, `LEDS_COUNT` and `FRAMERATE` accordingly. `IP` and `PORT` are of the controller's. Framerate is limited to 144 frames per second.

  ```
client.exe IP PORT LEDS_COUNT FRAMERATE
  ```

4. Double click the `run.bat` file to launch the client. It will autoconnect to the server once the server is online.


#### Server (e.g. Raspberry Pi Zero W)

1. Download the *server* binary for your operating system from [here](https://github.com/SHT/Ambilight/releases/latest/).
2. Install `tmux` using the following command:

  `sudo apt-get install tmux`

3. Create a `run.sh` file on the same folder as the binary. Make sure the file is marked as executable:

  `chmod +x server`

4. Paste the following on the `run.sh` file, replacing `AMBILIGHT_FOLDER`, `LEDS_COUNT`, `BRIGHTNESS`, `PIN` and `PORT` accordingly.  
  The arguments `PIN` and `PORT` are optional and default to `18` and `4197` respectively.

  ```
#!/bin/bash
tmux new-session -d -s ambilight 'cd /AMBILIGHT_FOLDER && ./server LEDS_COUNT BRIGHTNESS PIN PORT'
  ```

5. (optional) Start the server at boot: Edit the `/etc/rc.local` file, adding the following before the `exit 0` line, replacing `AMBILIGHT_FOLDER` with the folder where the ambilight server binary resides.

  ```
/AMBILIGHT_FOLDER/run.sh
  ```

6. Execute the run.sh file to start the Ambilight server (or reboot if you configured start at boot):

  ```
./run.sh
  ```

## Modes

```txt
A : Ambilight : Renders the streaming LED data that are provided from the client.
R : Rainbow   : Infinite loop of a gradient color shift animation.
```

## Under the hood

This library consists of 3 packages, a client (client.go), a server (server.go), and a utilities package (ambilight.go) that has functions that the client and server use to connect to one another and transmit/receive data, along with maintaining the state of the Ambilight service.

Once the client is online, it indefinitely attempts to connect to the server. When the server comes online, a TCP connection is established and the client starts capturing the screen and transmitting the LED color data to the server, using the *Ambilight* operation mode.

In this mode, the client is taking "screenshots", reading the border pixels, averaging them depending on the LEDs count, and sending the data to the server. The underlying [scrap](https://github.com/quadrupleslap/scrap) library captures raw pixel data from the client GPU's Backbuffer, which is a really more performant method than alternatives ([BitBlt](https://github.com/kbinani/screenshot), for example). There is no noticeable performance drop while the Ambilight mode is engaged.

When the socket connection is established, the client sends a message with this format:


```
Bytes (binary):

  - 1: MMMM MMMM * M is mode ascii character
  - 2: RRRR RRRR * R is red
  - 3: GGGG GGGG * G is green
  - 4: BBBB BBBB * B is blue
  - 5: RRRR RRRR * repeats for each additional LED
  - 6: GGGG GGGG
  - 7: BBBB BBBB
       ..
```

The first byte is the ASCII *mode* character.  
The rest of the bytes that follow MUST be `N * 3`, where `N` is the number of LEDs that will be controlled.  
If the strip has more or less LEDs the behavior is undefined.

The Rainbow mode is basically a moving color wheel gradient, meaning that all the LEDs have the next color in the chain, and they cycle all the available colors, from red to green to blue and back to red. It is implemented server-side, so no data needs to be sent after initiating it.


## Contributing
You are free and actively encouraged to contribute to this project by either contributing code, creating issues, reporting bugs, highlighting vulnerabilities, proposing improvements or helping maintain the documentation.

If you would like to submit code changes, create a new branch from the *develop* branch with the name of the feature you are implementing  and submit a pull request to the *develop* branch after you make your changes. Click [here](https://gist.github.com/Chaser324/ce0505fbed06b947d962#doing-your-work) for a how-to guide.

In case you want to submit a bug report, please add as many details as possible regarding how the error occured and include the steps required to reproduce it if that is possible. It will help a lot in testing, finding the cause and implementing fixes.

## Changelogs
Changelogs for each and every release can be found [here](https://github.com/SHT/Ambilight/releases).

## Copyright
Any reproductions of this project *must* include a link to this repository and the following copyright notice, along with the project's license.

© 2019 Tasos Papalyras - All Rights Reserved
