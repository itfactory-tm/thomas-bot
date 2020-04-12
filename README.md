Thomas Bot
==========

<img src="./images/logo.png" alt="Thomas Bot logo" width="150">

Thomas Bot is the friendly Discord bot! It hangs around in the official IT Factory Discord server.
It helps teachers doing their job and students also... sometimes... 

## Running locally

### Build manually
0. Make sure the [Go toolchain](https://golang.org/doc/install) is installed and working.
1. Make sure the required dependencies are installed:
    - On Ubuntu and other Debian-based distros, install the following packages: `libsox-dev libsdl2-dev portaudio19-dev libopusfile-dev libopus-dev curl`.
    - On Arch-based distros, install the following packages: `libsoxr sdl portaudio opusfile`.
    - On Windows, you're on your own for now :)
2. Build the project using `go build`.
3. Run the compiled `thomas-bot` binary:
    - On Unix-like systems, run `THOMASBOT_TOKEN={TOKEN} ./thomas-bot`, where `{TOKEN}` is your Discord bot's token.
    - On Windows, change your environment variables through your [system properties](https://docs.oracle.com/en/database/oracle/r-enterprise/1.5.1/oread/creating-and-modifying-environment-variables-on-windows.html). Then, run the bot.

### Docker
0. Make sure [Docker](https://docs.docker.com/get-started/) is installed and working.
1. Build the container using `docker build -t thomas-bot .`.
2. Run the image you've built using `docker run -it -e "THOMASBOT_TOKEN={TOKEN}" thomas-bot` where `{TOKEN}` is your Discord bot's token.
    - You can change the prefix by setting the `THOMASBOT_PREFIX` environment variable.
    - See main.go for more environment variables.
3. Interact with the bot using your chosen prefix, eg. `tm!love`.

## Credits
The cute robot is CC0 by Ann Hannes