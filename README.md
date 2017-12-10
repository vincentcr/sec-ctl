SecCtl

SecCtl is built out of 3 principal components:
 * `local`: an on-premise monitor daemon, collecting events and sending commands to the monitor process
 * `cloud`: a REST API to communicate with the daemon
 * `mock`: a mock TPI implementation for testing without access to physical device. Also useful for, eg, simluating alarms.

`local` and `cloud` are connected together with a web socket. `local` sends state changes to `cloud`, and `cloud`  can send commands to `local` through the socket.
