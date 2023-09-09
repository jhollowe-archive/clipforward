# TODO

* automatically detect the max number of characters that can be sent through the clipboard
* support multiple connections by round-robining them over the clipboard
* validate flags/config is the same/compatible between the server and client
* maybe use base64 encoding rather than raw bytes?
* maybe ACK packets since duplicate data will get ignored on the clipboard
* client should ping the server if it has not gotten any packet in a while (data or control)
