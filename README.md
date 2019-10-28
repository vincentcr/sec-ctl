## SecCtl

SecCtl is a cloud bridge to the Envisalink TCP/IP module for DSC and Honeywell alarm system panels. It monitors for and relays security events, as well as allowing sending commands to the panel through a REST API.

SecCtl is built out of 3 principal components:
 * `local`: an on-premise monitor daemon, collecting events from, and sending commands to, the Envisalink module;
 * `cloud`: a REST API to communicate with the daemon;
 * `mock`: a mock TPI implementation for testing without access to physical device. Also useful for, eg, simluating alarms.

`local` and `cloud` are connected together with a web socket. `local` sends state changes to `cloud`, and `cloud`  can send commands to `local` through the socket.
