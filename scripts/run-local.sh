#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"

export DYNAMODB_TABLE_NAME="linker_link"
export DYNAMODB_ENDPOINT="http://localhost:28000"
export IS_DEV=true
export DEV_PORT=8081

if [[ ! $# -eq 2 ]]; then
  echo "* Usage: $0 <compiled_filepath> <env_file>" && exit 1
fi

BUILD=$1
TEST_ENV=$2

. "$TEST_ENV"


_setup_dynamodb() {
  docker pull amazon/dynamodb-local:latest
  docker rm -f linker-dynamodb
  docker run --name linker-dynamodb \
    -d -p 28000:8000 amazon/dynamodb-local
}

_setup_dynamodb_table() {
  aws dynamodb create-table \
    --table-name "$DYNAMODB_TABLE_NAME" \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url "$DYNAMODB_ENDPOINT"

  aws dynamodb update-time-to-live --table-name "$DYNAMODB_TABLE_NAME" \
    --time-to-live-specification "Enabled=true, AttributeName=ttl" \
    --endpoint-url "$DYNAMODB_ENDPOINT"
}


echo "* Provisioning DynamoDB on local ..."

_setup_dynamodb &>/dev/null
_setup_dynamodb_table &>/dev/null

echo "* Running $BUILD"
$BUILD
