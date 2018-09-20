### flagger

*flagger* is the simplest feature flag service possible. Uses a JSON file to store flag state.

#### api

`GET /flags` returns a list of flags
`GET /flags/$FLAG_NAME/$ENVIRONMENT` returns the flag state for a particular environment. If environment is not found it will fall back to the `default` enviornment. If the `default` environment is not found it will fall back to a global default.

That's it.