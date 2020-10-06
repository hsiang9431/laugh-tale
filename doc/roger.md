# roger

The program that prepares encrypted files and directory structure for building encrypted images

## Execution flow

1. Read input files and prepare output folder structure
    a. Create `tar.gz` file from `payload` folder in `/input` volume mounted to its container
    b. Encrypt the `tar.gz` file with `AES256` encryption
    c. Copy all other required files: `poneglyph`, `poneglyph.sh` and `entrypoint.sh`
2. Second run to register encryption key ID with docker image ID. This step makes a http request
to `key management service` on endpoint `kozuki.service/v1/key/bind`

## Commands

Command | Flag | Description
--------|------|-------------
`encrypt` | \- | Encrypt `/input` volume of roger container to `/output` volume
\- | `--image` | The base image to use. Default to `ubuntu:latest`
\- | `--volume` | Allowed volumes\* in encrypted container, use this flag multiple times for adding more volumes
\- | `--key-server` | The URL of `key management service` (kozuki) to create keys
\- | `--verbose` | Print detailed log to console
`bind-image-id` | \- | Register docker image ID with encryption key ID
\- | `--key-id` | The encryption key ID for registration
\- | `--image-id` | The docker image ID for registration
\- | `--key-server` | The URL of `key management service` (kozuki) to create keys
\- | `--verbose` | Print detailed log to console


\* *volumes mounted right under root directory, `/`, triggers warning messages*

## Input

Volume | Directory name | File name | Description
-------|----------------|-----------|-------------
`/input` | \- | `entrypoint.sh` | The script for executing protected application in `payload`
`/input` | `payload` | \- | The directory that will be encrypted

## Output

Volume | File name | Description
-------|-----------|-------------
`/output` | `poneglyph` | The container implant binary acted as the entry point of this container
`/output` | `poneglyph.sh` | The utility script for escalating privilege and initiate second stage of `poneglyph`
`/output` | `entrypoint.sh` | The script for executing protected application in `payload`
`/output` | `enc_payload` | The encrypted `payload` tar.gz file
