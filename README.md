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

#### CLI

- `flagctl serve`: Start the HTTP API server on port 8082
- `flagctl set [FLAG_NAME] ([ENV_NAME]) [FLAG_TYPE] [FLAG_VALUE]`: Setup a flag
- `flagctl get [FLAG_NAME] ([ENV_NAME])`: Get a flags state

Note: ENV_NAME is optional but could be useful for customising flags based on environments

#### HTTP api usage

 - `GET /flags` returns a list of flags
 - `GET /flags/$FLAG_NAME/$ENVIRONMENT` returns the flag state for a particular environment. If environment is not found it will fall back to the `default` enviornment. If the `default` environment is not found it will fall back to a global default.

That's it.
