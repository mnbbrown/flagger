### flagger

*flagger* is the simplest feature flag service possible. Uses a JSON file to store flag state.

#### usage

- `kubectl serve`: Start
- `kubectl set [FLAG_NAME] ([ENV_NAME]) [FLAG_TYPE] [FLAG_VALUE]`: Setup a flag
- `kubectl get [FLAG_NAME] ({ENV_NAME})`: Get a flags value

#### api

`GET /flags` returns a list of flags
`GET /flags/$FLAG_NAME/$ENVIRONMENT` returns the flag state for a particular environment. If environment is not found it will fall back to the `default` enviornment. If the `default` environment is not found it will fall back to a global default.

That's it.
