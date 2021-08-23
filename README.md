# oyaki

[![CI](https://github.com/yano3/oyaki/actions/workflows/ci.yml/badge.svg)](https://github.com/yano3/oyaki/actions/workflows/ci.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/yano3/oyaki)](https://hub.docker.com/r/yano3/oyaki)

Dynamic image quality transformation proxy.

## Usage

### Docker

```
docker pull yano3/oyaki:latest
docker run -p 8080:8080 -e "OYAKI_ORIGIN_HOST=example.com" yano3/oyaki
```

## Configuration

Environment variables bellow are available.

- `OYAKI_ORIGIN_HOST`: Your origin host. Example: `example.com` (required)
- `OYAKI_ORIGIN_SCHEME`: Scheme to request to your origin. Default: `https`
- `OYAKI_QUALITY`: Image quality. Default: `90`

If using [Elastic APM](https://www.elastic.co/jp/apm/), please set the following environment variables.

- `ELASTIC_APM_SERVER_URL`: APM Server URL
- `ELASTIC_APM_SECRET_TOKEN`: Secret token for authentication. Optional.
- `ELASTIC_APM_SERVICE_NAME`: The name of service. Default: Use the executable name.

For other options, refer to [Elastic APM's Go Agent Reference](https://www.elastic.co/guide/en/apm/agent/go/current/configuration.html#configuration).
