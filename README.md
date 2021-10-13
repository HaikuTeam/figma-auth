# Figma Auth service for Haiku Animator

In order to use [Haiku Animator's](https://github.com/HaikuTeam/animator) Figma integration, a service must be running to perform OAuth2 token exchange.  Read more about OAuth with Figma and the Figma HTTP API here:  https://www.figma.com/developers/api#oauth2

This project is a minimal self-hosted HTTP server application written in Go that will perform the token exchange for you.

## Setup

0. Figure out the public URL and port where you will host this server.  This guide will use the made-up example http://animatorfigmaauth.na:8080/ . Using TLS/HTTPS is recommended and can be achieved by specifying the `TLS_*` environment variables in step 4. 
1. Log into Figma and register a new app at: https://www.figma.com/developers/apps .  The `name` and `logo` can be whatever you want.  The `website URL` must be the public URL from step 0 along with the path /v0/integrations/figma/token, e.g. `http://animatorfigmaauth.na:8080/v0/integrations/figma/token`
2. Add this callback URL to your registered app: `haiku://oauth/figma`
3. Save your app, then take note of the Figma-provided client ID and client secret — these values should be set to FIGMA_CLIENT_ID and FIGMA_CLIENT_SECRET in env (see .env.example)
4. Build and serve this application at the URL and port from step 0, having loaded the relevant `env` after copying `.env.example` into `.env` and filling in the relevant values. Once `.env` exists, you can use this one-liner to build, load env, and run the server: `go build && export $(grep -v '#.*' .env | xargs) && ./figma-auth`.  You can test that your server is accessible via `GET /v0/ping`, e.g. in your browser at http://animatorfigmaauth.na:8080/v0/ping .  If the server is accessible, you should see the response `pong`.
5. Specify that Haiku Animator should use your service for Figma auth — for now this requires building Animator from source `HAIKU_API=http://animatorfigmaauth.na:8080/ FIGMA_CLIENT_ID=get_this_from_figma yarn go` (from root of https://github.com/HaikuTeam/animator)

## License

This code is dual-licensed under Apache 2.0 and MIT.