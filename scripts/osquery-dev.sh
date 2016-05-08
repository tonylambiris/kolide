#!/usr/bin/env bash

# openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout private.key -out certificate.crt

export SECRET="kolidedev"
# SERVER="localhost:8000"
SERVER="localhost:8000"
CERT=./tmp/kolide.crt

echo $PWD

sudo osqueryd \
  --verbose \
  --pidfile /tmp/osquery.pid \
  --host_identifier uuid \
  --database_path /tmp/osquery.db \
  --config_plugin tls \
  --config_tls_endpoint /api/v1/osquery/config \
  --config_tls_refresh 10 \
  --config_tls_max_attempts 3 \
  --enroll_tls_endpoint /api/v1/osquery/enroll  \
  --enroll_secret_env SECRET \
  --disable_distributed=false \
  --distributed_plugin tls \
  --distributed_interval 10 \
  --distributed_tls_max_attempts 3 \
  --distributed_tls_read_endpoint /api/v1/osquery/read \
  --distributed_tls_write_endpoint /api/v1/osquery/write \
  --tls_dump true \
  --logger_path /tmp/ \
  --logger_plugin tls \
  --logger_tls_endpoint /api/v1/osquery/log \
  --logger_tls_period 5 \
  --tls_hostname $SERVER \
  --tls_server_certs $CERT \
  --pack_delimiter /
