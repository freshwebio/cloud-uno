#!/bin/bash

POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -t|--tests)
    TEST_TYPES="$2"
    shift # past argument
    shift # past value
    ;;
    -r|--run)
    RUN_LOCALLY=yes
    shift # past argument
    ;;
    -s|--stop)
    STOP=yes
    shift # past argument
    ;;
    -c|--clean)
    CLEAN_UP_LOCAL_RUN=yes
    shift
    ;;
    --server-in-docker)
    RUN_SERVER_IN_DOCKER=yes
    shift # past argument
    ;;
    -h|--help)
    HELP=yes
    shift # past argument
    ;;
    *)    # unknown option
    POSITIONAL+=("$1") # save it in an array for later
    shift # past argument
    ;;
esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

function help {
  cat << EOF
Service tools
Provides tooling to run the server and host agent for development, run tests and deal with set up and teardown
for integration tests.

To run the server (Direct on host) and host agent locally for development with live reload:
bash scripts/service-tools.sh --run

To run the server (in Docker) and host agent locally for development with live reload:
bash scripts/service-tools.sh --run --server-in-docker

To stop the server and host agent running locally:
bash scripts/service-tools.sh --stop

To run the server locally and clean up docker containers on exit:
bash scripts/service-tools.sh --run --clean

To run the server locally (in Docker) and clean up docker containers on exit:
bash scripts/service-tools.sh --run --clean --server-in-docker

To run integration and unit tests:
bash scripts/service-tools.sh --tests unit,integration

To run only unit tests:
bash scripts/service-tools.sh --tests unit

To run only integration tests:
bash scripts/service-tools.sh --tests integration
EOF
}

if [ -n "$HELP" ]; then
  help
  exit 0
fi

function finish {
  # Make every test run repeatable by ensuring we tear down the dependencies.
  # Optionally, clean up when running locally if desired.
  if [ -n "$TEST_TYPES" ] && [ "$TEST_TYPES" != "unit" ]; then
    compose_file="docker-compose.tests.yml"
  elif [ -n "$CLEAN_UP_LOCAL_RUN" ]; then
    compose_file="docker-compose.local.yml"
  fi

  if ([ -n "$TEST_TYPES" ] && [ "$TEST_TYPES" != "unit" ]) || [ -n "$CLEAN_UP_LOCAL_RUN" ]; then
    echo "Tearing down docker dependencies ..."
    docker-compose --file $compose_file stop
    docker-compose --file $compose_file rm -v -f
  fi

  if [ -z "$TEST_TYPES" ]; then
    echo "Tearing down the local server docker container if needed ..."
    # If the server is running in docker then let's make sure we stop it.
    docker-compose --file docker-compose.server-local.yml stop
    docker-compose --file docker-compose.server-local.yml rm -v -f

    # Kill the Go server if it is running directly on
    # the host.
    server_pid=$(lsof -t -i :5988 -s TCP:LISTEN)
    if [ -n "$server_pid" ]; then
      echo "Killing the local server process running on the host ..."
      # Ensure we kill the air live reload process,
      # this will also kill the child process running the server.
      air_pid=$(ps -o ppid= -p $server_pid)
      if [ -n "$air_pid" ]; then
        kill -9 $air_pid
      fi
    fi
  fi

  # Kill the host agent along with
  # the air process that controls it if it is running.
  # pgrep returns a 1 exit code if no matches are found which will fail the step
  # in CI environments so we need to ensure it returns a 0 exit code.
  hostagent_air_pid=$(pgrep -f "air -c \.air\.hostagent\.toml" || true)
  if [ -n "$hostagent_air_pid" ]; then
    echo "Killing the host agent ..."
    kill -9 $hostagent_air_pid
  fi
}
trap finish EXIT

if [  -n "$STOP" ]; then
  finish
  # If stop is called explicitly then let's exit early.
  exit 0
fi

if ([ -n "$TEST_TYPES" ] && [ "$TEST_TYPES" != "unit" ]) || [ -n "$RUN_LOCALLY" ]; then
  echo "Bringing up dependency services (Databases, caches, emulators etc.)"

  docker_compose_file="docker-compose.local.yml"
  if [ -n "$TEST_TYPES" ]; then
    docker_compose_file="docker-compose.tests.yml"
  fi

  docker-compose --file "$docker_compose_file" up -d

  echo "Waiting a few seconds to allow time for dependencies to be available ..."
  sleep 10s
fi

if [ -n "$TEST_TYPES" ]; then
  source .testrc

  set -e
  echo "" > coverage.txt

  go test -timeout 30000ms -tags "$TEST_TYPES" -race -coverprofile=coverage.txt -covermode=atomic ./...

  if [ -n "$GITHUB_ACTION" ]; then
    # We are in a CI environment so run tests again to generate JSON report.
    go test -timeout 30000ms -json -tags "$TEST_TYPES" ./... > report.json
  fi

else
  source .localrc

  # Regardless if the server is running in Docker or not, we need to build the client on the host machine
  # as the docker image used to run the development Docker server does not have nodejs and npm installed.
  pushd client
  echo "Building client (You will need to re-build the client with \"yarn build\" to reflect any changes you make) ..."
  yarn build
  popd

  echo "Starting the live-reload Cloud::1 host agent as root ..."
  # The host agent should be running as root user as it manipulates the hosts file.
  sudo bash -c "source .localrc && air -c .air.hostagent.toml 2>&1 > air-hostagent.log &"

  echo "Waiting 5 seconds for host agent to start connection on socket ..."
  sleep 5s

  echo "Running the live-reload Cloud::1 API server ..."
  if [ -n "$RUN_SERVER_IN_DOCKER" ]; then
    docker-compose -f docker-compose.server-local.yml up -d
  else
    air -c .air.server.toml 2>&1 > air-server.log &
  fi

  echo "Server and host agent are running, use CTRL+C to exit ..."
  # idle waiting for abort from user
  read -r -d '' _ </dev/tty
fi