#!/usr/bin/env bash

# openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout private.key -out certificate.crt

CERT=${CERT:-./tmp/kolide.crt}
KEY=${KEY:-./tmp/kolide.key}

CN="$(openssl x509 -in $CERT -text -noout 2>/dev/null | \
	awk '$1 ~ /Subject:/ {print $6}' | cut -d '=' -f 2)"
SERVER="${SERVER:-${CN}:8000}"

export SECRET=${SECRET=kolidedev}

echo
echo "$PWD (SERVER=$SERVER CERT=$CERT SECRET=$SECRET)"
echo

sudo -E osqueryd \
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
  --tls_dump=true \
  --logger_path /tmp/ \
  --logger_plugin tls \
  --logger_tls_endpoint /api/v1/osquery/log \
  --logger_tls_period 5 \
  --tls_hostname $SERVER \
  --tls_server_certs $CERT \
  --tls_client_cert $CERT \
  --tls_client_key $KEY \
  --pack_delimiter /
