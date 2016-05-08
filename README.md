[![Build Status](http://komanda.io:8080/api/badges/mephux/kolide/status.svg)](http://komanda.io:8080/mephux/kolide)

# Kolide

  Kolide is an agentless osquery web interface and remote api server. Kolide was designed to be extremely portable 
  and performant while also keeping a simple codebase for others to contribute to. I have a lot planned
  for Kolide so check back often! :)

# Kolide Rest API & OSquery Remote API

  Kolide has a rest api that uses jwt (https://jwt.io/) for authentication. The only exception
  to this is the osquery remote api below.

  [osquery remote api](https://osquery.readthedocs.org/en/stable/deployment/remote/#remote-server-api):

  method | url | osquery configuration cli flag
  -------|-----|-------------------------------
  POST   | /api/v1/osquery/enroll | `--enroll_tls_endpoint`
  POST   | /api/v1/osquery/config | `--config_tls_endpoint`
  POST   | /api/v1/osquery/log    | `--logger_tls_endpoint`
  POST   | /api/v1/osquery/read   | `--distributed_tls_read_endpoint`
  POST   | /api/v1/osquery/write  | `--distributed_tls_write_endpoint`

# TODO

  What else should I add? Add a ticket and let me know!

  - [X] Database auth (salt/hash)
  - [ ] LDAP support
  - [X] Full OSquery Remote API Support
  - [X] Read/Write UI
  - [X] Basic saved queries and loading
  - [X] Websockets (live node updates)
  - [X] CI 
  - [ ] Websockets (osquery logs)
  - [ ] OSquery Config & Pack
  - [ ] Node Editing
  - [ ] Scheduled query UI
  - [ ] Rules engine for alerting (slack, pagerduty etc..)
  - [ ] Auto build and publish deb/rpms using packagecloud.io
  - [ ] Write tests for reasonable things

# Development

  The easiest way to start writing code is to use docker/docker-compose.

  * `make up` will run docker-compose and bootstrap the deps
  * `make down` will spin down and remove all deps

## osqueryd

  ```bash
  #!/usr/bin/env bash

  export SECRET="kolidedev"
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
  ```

# Contributing

When contributing to this repository, please first discuss the change you wish to make via issue,
email, or any other method with the owners of this repository before making a change. 

Please note we have a code of conduct, please follow it in all your interactions with the project.

## Pull Request Process

1. Ensure any install or build dependencies are removed before the end of the layer when doing a 
   build.
2. Update the README.md with details of changes to the interface, this includes new environment 
   variables, exposed ports, useful file locations and container parameters.
3. Increase the version numbers in any examples files and the README.md to the new version that this
   Pull Request would represent. The versioning scheme we use is [SemVer](http://semver.org/).
4. You may merge the Pull Request in once you have the sign-off of two other developers, or if you 
   do not have permission to do that, you may request the second reviewer to merge it for you.

# Self-Promotion

Like kolide? Follow the repository on
[GitHub](https://github.com/mephux/kolide) and if
you would like to stalk me, follow [mephux](http://dweb.io/) on
[Twitter](http://twitter.com/mephux) and
[GitHub](https://github.com/mephux).
