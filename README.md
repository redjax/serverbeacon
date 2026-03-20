# Server Beacon

Simple Go app that listens on multiple interfaces for external hosts to check server's online status.

## Progress

- [ ] HTTP server/REST API
  - [x] Simple HTTP server
  - [x] `/ping` endpoint (return "pong")
  - [ ] `/health` endpoint
    - [x] Return JSON with health status and timestamp
    - [ ] An authenticated `/health` endpoint that returns stats about the underlying host
- [ ] SSH server
  - [ ] Allow SSH connections, immediately terminate
  - [ ] Only allow SSH keys, no user/password auth
- [ ] RPC message
  - [ ] Create a client to 'ping' the server
- [ ] Docker container
