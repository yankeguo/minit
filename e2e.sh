#!/bin/bash

set -eu

echo "Running end-to-end tests"

GOOS=linux GOARCH=amd64 go build -o minit

docker build -t minit:e2e -f Dockerfile.e2e .

DIR_TEMP="$(mktemp -d)"

mkdir -p "$DIR_TEMP"

echo "1. General units loading"

DIR_CASE="$DIR_TEMP/case-1"

mkdir -p "$DIR_CASE/minit.d"

cat <<-EOF >"$DIR_CASE/minit.d/unit1.yaml"
kind: once
name: file-once
group: file
count: 2
critical: true
dir: /tmp
shell: /bin/sh
env:
  a: b
charset: gbk
command:
  - echo -n 'xOO6wwo=' | base64 -d \
  - | cat
  - exit 2
success_codes: [0, 1, 2]
EOF

docker run -ti --rm -e MINIT_QUICK_EXIT=true -v "$DIR_CASE/minit.d:/etc/minit.d" minit:e2e
