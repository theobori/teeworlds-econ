#!/usr/bin/env bash

set -e

readonly ECON_DEBUG=${ECON_DEBUG:-"1"}
readonly ECON_PORT=${ECON_PORT:-"7000"}
readonly ECON_PASSWORD=${ECON_PASSWORD:-"hello_world"}

readonly CONTAINER_NAME="Teeworlds"

# Create the override configuration file
cat << EOF > myServerconfig.cfg
ec_port ${ECON_PORT}
ec_password "${ECON_PASSWORD}"
ec_output_level 2
EOF

# Create a container
docker run --name ${CONTAINER_NAME} -d \
  -p ${ECON_PORT}:${ECON_PORT} \
  -p 8303:8303 \
  -p 8303:8303/udp \
  -v $PWD/myServerconfig.cfg:/serverdata/serverfiles/DDNet/myServerconfig.cfg \
  --env 'GAME_CONFIG=autoexec.cfg' \
  --env 'UID=1000' \
  --env 'GID=1000' \
  ich777/ddnetserver:latest
