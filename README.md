# Dev

1. Install [devbox](https://github.com/jetpack-io/devbox#installing-devbox)
2. Visit https://twitchtokengenerator.com
3. Set ClientId (TWITCH_CLIENTID) and ClientSecret (TWITCH_CLIENTSECRET) of twitch from doppler settings
4. Generate token with all scopes
5. Set tokens to the .env file:
   1. access token as `BOT_ACCESS_TOKEN`
   2. refresh token as `BOT_REFRESH_TOKEN`
   3. `DOPPLER_TOKEN` from your doppler envirement
6. Enter devboxshell:
   ```bash
   devbox shell
   ```
7. Start all services:
   ```bash
   task dev
   ```

### IMPORTANT

For exiting envirement, you should use `CTRL+D` hostkey TWICE!
