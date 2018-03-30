#!/bin/bash

set -e

g_cmd_clearchainctl=${CLEARCHAINCTL:-./clearchainctl}
g_cmd_clearchaind=${CLEARCHAIND:-./clearchaind}
g_chadmin_key_name=chadmin
g_chadmin_key_pub=''
g_chain_id=''

f_create_key_return_pub() {
  local l_key_name l_key_pub

  l_key_name="$1"
  "${g_cmd_clearchainctl}" keys add "${l_key_name}" 1>/dev/null 2>&1 <<EOF
password
password
EOF
  l_key_pub="`${g_cmd_clearchainctl} pub2hex ${l_key_name}`"
  printf "%s" "${l_key_pub}"
  return 0
}

f_init_chain_with_pub() {
  local l_pub

  l_pub="$1"
  "${g_cmd_clearchaind}" init "${l_pub}"
}

f_cleanup() {
  rm -vrf .clearchainctl .clearchaind
}

f_start_clearchaind() {
  "${g_cmd_clearchaind}" start &>clearchaind.log
  printf "%s" "`grep chain_id .clearchaind/config/genesis.json | sed -n 's/^.*: \"\(.*\)\".*/\1/p'`"
  return 0
}

f_cleanup
g_chadmin_key_pub="`f_create_key_return_pub ${g_chadmin_key_name}`"
f_init_chain_with_pub "${g_chadmin_key_pub}"
echo "Clearing House admin key: ${g_chadmin_key_name} - pub: ${g_chadmin_key_pub}"
g_chain_id=`f_start_clearchaind`
