# Dev

1. Visit https://twitchtokengenerator.com
2. Set ClientId (TWITCH_CLIENTID) and ClientSecret (TWITCH_CLIENTSECRET) of twitch from doppler settings
3. Generate token with all scopes
4. Set tokens to the .env file:
   1. access token as `BOT_ACCESS_TOKEN`
   2. refresh token as `BOT_REFRESH_TOKEN`
   3. doppler token from (https://dashboard.doppler.com/workplace/fa99234a9a83eb6bb7e8/projects/tsuwari/configs/dev/access) as `DOPPLER_TOKEN`
5. For build use `make build-dev`
6. For start use `make dev`
