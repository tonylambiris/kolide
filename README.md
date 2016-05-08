[![Build Status](http://komanda.io:8080/api/badges/mephux/kolide/status.svg)](http://komanda.io:8080/mephux/kolide)

# Kolide

  Kolide is an agentless osquery web interface and remote api server. Kolide was designed to be extremely portable 
  and performant while also keeping a simple codebase for others to contribute to. I have a lot planned
  for Kolide so check back often! :)

# Remote API endpoints

  [osquery remote api](https://osquery.readthedocs.org/en/stable/deployment/remote/#remote-server-api):

  method | url | osquery configuration cli flag
  -------|-----|-------------------------------
  POST   | /api/v1/osquery/enroll | `--enroll_tls_endpoint`
  POST   | /api/v1/osquery/config | `--config_tls_endpoint`
  POST   | /api/v1/osquery/log    | `--logger_tls_endpoint`
  POST   | /api/v1/osquery/read   | `--distributed_tls_read_endpoint`
  POST   | /api/v1/osquery/write  | `--distributed_tls_write_endpoint`

# osqueryd

    ~~~
    #!/usr/bin/env bash
    
    export ENROLL_SECRET=secret
    
    osqueryd \
       --pidfile /tmp/osquery.pid \
       --host_identifier uuid \
       --database_path /tmp/osquery.db \
       --config_plugin tls \
       --config_tls_endpoint /config \
       --config_tls_refresh 10 \
       --config_tls_max_attempts 3 \
       --enroll_tls_endpoint /enroll  \
       --enroll_secret_env ENROLL_SECRET \
       --disable_distributed=false \
       --distributed_plugin tls \
       --distributed_interval 10 \
       --distributed_tls_max_attempts 3 \
       --distributed_tls_read_endpoint /distributed/read \
       --distributed_tls_write_endpoint /distributed/write \
       --tls_dump true \
       --logger_path /tmp/ \
       --logger_plugin tls \
       --logger_tls_endpoint /log \
       --logger_tls_period 5 \
       --tls_hostname localhost:5000 \
       --tls_server_certs ./certificate.crt \
       --log_result_events=false \
       --pack_delimiter /
    ~~~

# Development

  The easiest way to start writing code is to use docker/docker-compose.

  * `make up` will run docker-compose and bootstrap the deps
  * `make down` will spin down and remove all deps
