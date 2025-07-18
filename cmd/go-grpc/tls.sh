
export SERVER_ADDRESS=":50051"
export SERVER_USE_TLS="true"
export SERVER_CERT_FILE="/workspaces/playground/certs/server.cert.pem"
export SERVER_KEY_FILE="/workspaces/playground/certs/server.key.pem"
export SERVER_CA_FILE="/workspaces/playground/certs/ca.cert.pem"
export SERVER_SERVER_NAME=""
export SERVER_CLIENT_CERT_AUTH="true"

export SERVER_METRICS_ENABLED="true"
export SERVER_METRICS_ADDRESS=":9092"
export SERVER_METRICS_AUTH_USERNAME="user"
export SERVER_METRICS_AUTH_PASSWORD="password"

export CLIENT_ADDRESS="dns:///localhost:50051"
export CLIENT_USE_TLS="true"
export CLIENT_CERT_FILE="/workspaces/playground/certs/client.cert.pem"
export CLIENT_KEY_FILE="/workspaces/playground/certs/client.key.pem"
export CLIENT_CA_FILE="/workspaces/playground/certs/ca.cert.pem"
export CLIENT_SERVER_NAME=""
export CLIENT_KEEP_ALIVE="7s"
