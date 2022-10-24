# Dev

1. Install [devbox](https://github.com/jetpack-io/devbox#installing-devbox)
2. Visit https://twitchtokengenerator.com
3. Set `TWITCH_CLIENTID` and `TWITCH_CLIENTSECRET` of twitch from doppler settings
4. Generate token with all scopes
5. Set tokens to the .env file:
   1. access token as `BOT_ACCESS_TOKEN`
   2. refresh token as `BOT_REFRESH_TOKEN`
   3. `DOPPLER_TOKEN` from your doppler envirement
6. Enter `devbox shell`:
   ```bash
   devbox shell
   ```
7. Use `doppler setup` and select your dev env.
8. Start infrostructure:
   ```bash
   make up-dev
   ```
9. Start services:
   ```bash
   task dev
   ```

### IMPORTANT

For exiting envirement, you can use `CTRL+D` or type `exit` in terminal.
