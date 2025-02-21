# Go-Node gRPC

Helloworld example to learn gRPC with Node.js + Typescript on the client and Golang on the server, Server authentication with SSL/TLS.

Protoc plugin for TypeScript: [ts-proto](https://www.npmjs.com/package/ts-proto)

## Try it

- Clone the repository:

  ```bash
  git clone git@github.com:AnhBigBrother/go-node-gRPC.git
  ```

- Gen proto, certificate & install dependencies:

  ```bash
  make setup
  ```

- Run server:

  ```bash
  make server
  ```

- Run client(on another terminal):

  ```bash
  cd client
  npx ts-node main.ts [cmd] [your_message_to_server] 
  # Avaiable cmd: say_hello, stream_reply, stream_request, bidirection_stream
  # Example: npx ts-node main.ts --cmd=stream_request 'Hello from me'
  ```
