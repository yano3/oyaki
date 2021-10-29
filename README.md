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
