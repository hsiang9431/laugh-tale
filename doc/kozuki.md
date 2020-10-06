# Kozuki

## Functionalities

The service image deploy one of the following service according to environment variable `ENABLE_CRUD_SERVER`

- Serve APIs for `roger` and `ohara` with following functions on REST endpoints
- Serve APIs for general CRUD operations for directly managing keys in database

## Environment variables

Key | Type | Description
----|------|-------------
ENABLE_CRUD_SERVER | `bool` | If this variable is set or not empty, the server will serve CRUD APIs
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
DB_HOST | `string` | Database host
DB_PORT | `string` | Database port
DB_NAME | `string` | Database name in database server
DB_USERNAME | `string` | Database username
DB_PASSWORD | `string` | Database password
DB_TLS_MODE | `string` | Database TLS mode
DB_TLS_CERT | `string` | File path to database server TLS certificate


## `REST` endpoints

`Key` reference to `pkg/kozuki/types/Key` structure

### Option without `ENABLE_CRUD_SERVER` set

API | Method | Post form | Return value | Description
----|--------|-----------|--------------|-------------
`/v1/key/create` | `POST` | \- | `Key` structure in `json` with `key_id` set | Create a new encryption key in database
`/v1/key/bind` | `POST` | \- | `Key` structure in `json` with `key_id` and `image_id` set | The API for registering image ID to key by its ID
\- | \- | `image_id` | \- | The image ID for registration
\- | \- | `key_id` | \- | The key ID for registration
`/v1/key/{image_id}` | `GET` | \- | `Key` structure in `json` with required keys | Retrieve decryption keys registered to `image_id`


### Option with `ENABLE_CRUD_SERVER` set

All of the APIs in this category returns Key structure in `json`
containing effected `pkg/kozuki/types/Key` structure

API | Method | Body or Parameter | Description
----|--------|-------------------|-------------
`/key` | `POST` | Key structure in `json` | Create operation on database
`/key` | `GET` | `image_id` or `key_id` | Retrieve operation on database
`/key` | `PATCH` | Key structure in `json` | Update operation on database
`/key` | `DELETE` | Key structure in `json` | Delete operation on database
