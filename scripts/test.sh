#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"

_setup_dynamodb() {
  docker pull amazon/dynamodb-local:latest
  docker rm -f linker-dynamodb
  docker run --name linker-dynamodb \
    -d -p 28000:8000 amazon/dynamodb-local
}

_init() {
  aws dynamodb create-table \
    --table-name "$DYNAMODB_TABLE_NAME" \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url http://localhost:28000

  aws dynamodb update-time-to-live --table-name "$DYNAMODB_TABLE_NAME" \
    --time-to-live-specification "Enabled=true, AttributeName=ttl" \
    --endpoint-url http://localhost:28000
}


echo "* Provisioning DynamoDB on local ..."
export DYNAMODB_TABLE_NAME="linker_link"
_setup_dynamodb &>/dev/null
_init &>/dev/null || (_term; exit 1)

TEST_ENV="$SCRIPT_DIR/../test/env"
if [[ ! -e "$TEST_ENV" ]]; then
  echo "* [ERROR] Cannot find $TEST_ENV" && exit 1
fi

. "$TEST_ENV"

BUILD=$1

export LINKER_IS_DEV=true
export LINKER_DEV_PORT=8081
echo "Run $BUILD"
$BUILD
