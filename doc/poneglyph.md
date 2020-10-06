# poneglyph

## Execution flow

1. First stage run with `poneglyph start`. Which will connect to `key retriever service` on `/v1/key` end point
for obtaining decryption key. If succeed, it calls into `poneglyph.sh` for executing next stage.
2. Second stage executes `poneglyph run` with payload decryption key. This stage will decrypt
`enc_payload` back to `payload` directory and execute `entrypoint.sh`.

## Commands

Command | Flag | Description
--------|------|-------------
`start` | \- | Fetch keys from `key retriever service` found within the same Kubernetes cluster of current container
`run` | \- | Decrypt payload and executes `entrypoint.sh`
\- | `--dec-key-b64` | The decryption key for payload

## Script for executing 2nd stage

```shell
echo $1 | sudo -S -- sh -c "chmod o-rwx /*; /secrete-container/bob run --dec-key-b64 $2"
```
