Thomas Bot
==========

<img src="./images/logo.png" alt="Thomas Bot logo" width="150">

Thomas Bot is the friendly Discord bot! It hangs around in the official IT Factory Discord server.
It helps teachers doing their job and students also... sometimes... 

## Running locally

### Docker
1. Build the container using `docker build -t thomas-bot .`.
2. Run the image you've built using `docker run -it -e "THOMASBOT_TOKEN={TOKEN}" thomas-bot` where `{TOKEN}` is your Discord bot's token.
    - You can change the prefix by setting the `THOMASBOT_PREFIX` environment variable.
    - See main.go for more environment variables.
3. Interact with the bot using your chosen prefix, eg. `tm!love`.

## Credits
The cute robot is CC0 by Ann Hannes