mkdir -p certificates
cd certificates

# generate self signed ca.key and ca.crt as root ca
openssl req -new -text -passout pass:abcd -subj /CN=localhost -out ca.req
openssl rsa -in privkey.pem -passin pass:abcd -out ca.key
openssl req -x509 -in ca.req -text -key ca.key -out ca.crt

# generate a key certificate pair for roger
openssl req -new -text -passout pass:abcd -subj /CN=localhost -out roger.req
openssl rsa -in privkey.pem -passin pass:abcd -out roger.key
openssl req -x509 -in roger.req -text -key ca.key -out roger.crt

# generate a key certificate pair for kozuki
openssl req -new -text -passout pass:abcd -subj /CN=localhost -out kozuki.req
openssl rsa -in privkey.pem -passin pass:abcd -out kozuki.key
openssl req -x509 -in kozuki.req -text -key ca.key -out kozuki.crt
