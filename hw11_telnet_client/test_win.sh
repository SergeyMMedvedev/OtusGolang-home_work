#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | ncat -l localhost 4242 >"F:\dev\Otus\OtusGolang-home_work\hw11_telnet_client\tmp\nc.out" &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=5s localhost 4242 >"F:\dev\Otus\OtusGolang-home_work\hw11_telnet_client\tmp\telnet.out" &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}" && exit 1)
}

expected_nc_out='I
am
TELNET client'
fileEquals "F:\dev\Otus\OtusGolang-home_work\hw11_telnet_client\tmp\nc.out" "${expected_nc_out}"

expected_telnet_out='Hello
From
NC'
fileEquals "F:\dev\Otus\OtusGolang-home_work\hw11_telnet_client\tmp\telnet.out" "${expected_telnet_out}"

rm -f go-telnet
echo "PASS"


ncat -l localhost 4242 >/dev/null