# Goreprox
A small reverse port forwarding proxy written in go. 

*Development state: Pre-Alpha*

Lets you expose a local port via another machine you control. Just like what port forwarding services do, but all under your control.

# Requirements
Description of what the components are supposed to do.

## General 
- The code should be clean as hell A+++
- Performance is not the main goal but is important. The proxying should not add more than 100ms to the request roundtrip when tested localy.
- TCP only. HTTP port forwarding is the main goal but all TCP stuff should work.
- Must support multiple connections in parallel. A max amount may be configurable.

## Goreprox Server
- Listens on port 8080 for client requests. Port must be configurable.
- Listens on port 9887 for one (optionally multiple) providers. Port must be configurable.

- For all accepted connections on port 8080 forward the received packages to the provider. 
- All data comming from the provider must be forwarded to the client connection it is for.

- If an client connection is closed that must be communicated to the provider.
- If the connection to the provider is lost the server must start listening for another provider connection. Client connections must be kept intact for at least 1 minute.
- If the provider delivers the information that one of his target server connections was closed the client connection used that uses this target connection must be closed.
- (optional) clients are authenticated (eg. via custom header X-Goreprox-Auth)

## Provider
- Connects to the given server and port on startup. 
- Opens a connection to the target server & port.
- All data received from the server must be forwarded to the target Server
- Data from the same inbound connection on the server must be forwarded to the target server using the same connection.
- If the target server closes the connection that must be communicated to the Goreprox server
- If the goreprox server provides the information that one of his inbound connections was closed the target server connection used to forward the packages for that client connection must be closed.
- When the connection to the Goreprox server is lost the provider must reconnect and keep the target server connections intact for at least 1 minute


