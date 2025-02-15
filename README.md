# Go-Node gRPC

Helloworld example to learn gRPC with Node.js + Typescript on the client and Golang on the server.

Protoc plugin for TypeScript: [ts-proto](https://www.npmjs.com/package/ts-proto)

## Try it

- Clone the repository:

  ```bash
  git clone git@github.com:AnhBigBrother/go-node-gRPC.git
  ```

- Setup server:

  ```bash
  cd server
  sh gen_proto.sh
  go mod tidy
  go run main.go
  ```

- Setup client:

  ```bash
  cd client
  npm install
  sh gen_proto.sh
  ```

  ```bash
  npx ts-node main.ts [cmd] [your_message_to_server] 
  # Avaiable cmd: say_hello, stream_reply, stream_request, bidirection_stream
  # Example: npx ts-node main.ts --cmd=stream_request 'Hello from me'
  ```
