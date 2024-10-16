# Gommo
    simple cli-based mmo server from scratch in go

## Client
send "gommo" to authenticate
on authentication, request packet containing server for client to hold / render
- update server on move
- update on disconnect/socket close
needs map renderer

## Server
wait for "gommo" requests
- multithreaded handling
- want list of these connections eventually
send map packet upon request
holds map that can be sent out

## Shared
contains important stuff both programs need to know

packet syntax:
"<packetLength>\n<gommo>\n<packetType>\n<data/nothing>"
