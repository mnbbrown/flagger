## flagger

*flagger* is a simple feature flag microservice backed by redis.

### Installation

#### docker
```bash
docker run -d -p 8082:8082 mnbbrown/flagger serve
```

#### bash
```bash
wget https:///
cp flagctl /usr/local/bin/flagctl
sudo chmod +x /usr/local/bin/flagctl
```

### Usage

#### via CLI

- `flagctl serve`: Start the HTTP API server on port 8082
- `flagctl set [FLAG_NAME] ([ENV_NAME]) [FLAG_TYPE] [FLAG_VALUE]`: Setup a flag
- `flagctl get [FLAG_NAME] ([ENV_NAME])`: Get a flags state

Note: ENV_NAME is optional but could be useful for customising flags based on environments

#### via HTTP API

 - `GET /flags` returns a list of flags
 - `GET /flags/$FLAG_NAME/$ENVIRONMENT` returns the flag state for a particular environment. If environment is not found it will fall back to the `default` enviornment. If the `default` environment is not found it will fall back to a global default.

#### via Libraries

 - [Go](client)
 - [Javascript](https://github.com/mnbbrown/flagger-js-client)

#### Reference

There are two differnt types of flags `BOOL` and `PERCENT`.

- `BOOL` returns the same value all the time (i.e either true or false).
- `PERCENT` returns true $PERCENT% of the time. Useful for things like doing blue green deployments or tracing
