### flagger

*flagger* is the simplest feature flag service possible. Uses a JSON file to store flag state.

#### cli usage

- `kubectl serve`: Start
- `kubectl set [FLAG_NAME] ([ENV_NAME]) [FLAG_TYPE] [FLAG_VALUE]`: Setup a flag
- `kubectl get [FLAG_NAME] ([ENV_NAME])`: Get a flags state

Note: ENV_NAME is optional but could be useful for customising flags based on environments

#### HTTP api usage

 - `GET /flags` returns a list of flags
 - `GET /flags/$FLAG_NAME/$ENVIRONMENT` returns the flag state for a particular environment. If environment is not found it will fall back to the `default` enviornment. If the `default` environment is not found it will fall back to a global default.

That's it.
