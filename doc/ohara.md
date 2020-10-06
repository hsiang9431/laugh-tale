# Ohara

## Functionalities

1. Serve REST endpoints for `poneglyph` to retrieve decryption keys
2. Retrieve `deployment` details from `Kubernetes server API` and access `/v1/key/{image_id}` endpoint on `key management service` with docker image ID extracted from `deployment` details

## Environment variables

Key | Type | Description
----|------|-------------
K8S_API_SERVER | `string` | Kubernetes API server URL, should be set by Kebernetes controller
TLS_CERT_FILE | `string` | The TLS certificate file for serving the server
TLS_KEY_FILE | `string` | The TLS key file for serving the server
SERVER_LOG_FILE | `string` | The output log file for server. If pre-pend with `debug:`, server enables debug log
HTTP_LOG_FILE | `string` | The combined log file for http server. Such log will not be generated if this variable is not set
SERVER_PORT | `int` | The port to run server
READ_TIMEOUT | `time.Duration` | Server read timeout. The value is parsed with `time.ParseDuration`
READ_HEADER_TIMEOUT | `time.Duration` | Server read header timeout. The value is parsed with `time.ParseDuration`
WRITE_TIMEOUT | `time.Duration` | Server write timeout. The value is parsed with `time.ParseDuration`
IDLE_TIMEOUT | `time.Duration` | Server idle timeout. The value is parsed with `time.ParseDuration`
MAX_HEADER_BYTES | `int` | Server max header bytes


## `REST` endpoints

API | Method | Return value | Description
----|--------|--------------|-------------
`/v1/key` | `GET` | `pkg/ohara/types/EncKey` structure in `json` | Contains two keys required to decrypt payload
