import parseArgs from "minimist"
import { GreeterClient, HelloRequest, HelloReply } from "./proto/helloworld"
import * as grpc from "@grpc/grpc-js"

async function wait(time: number) {
  return new Promise((resolve) => setTimeout(resolve, time))
}

function runSayHello(client: GreeterClient, request: HelloRequest) {
  client.sayHello(request, function (err: grpc.ServiceError | null, response: HelloReply) {
    console.log("Received from server:", response.message)
  })
}

function runStreamRequest(client: GreeterClient, request: HelloRequest) {
  const streamRequest = client.sayHelloStreamRequest(
    (err: grpc.ServiceError | null, response: HelloReply) => {
      if (err) {
        console.log(err)
        return
      }
      console.log("Received from server:", response.message)
    }
  )

  const streamWrite = async (num: number) => {
    streamRequest.write({
      message: `${request.message} no.${num}`,
    })
  }
  const streamConcurent = async () => {
    for (let i = 0; i < 10; i++) {
      await streamWrite(i + 1)
      await wait(500)
    }
    streamRequest.end()
  }

  streamConcurent()
}

function runStreamReply(client: GreeterClient, request: HelloRequest) {
  const streamReply = client.sayHelloStreamReply(
    request,
    new grpc.Metadata({ cacheableRequest: false })
  )
  streamReply.on("data", (data) => {
    console.log("Received from server:", data.message)
  })
  streamReply.on("end", () => {
    console.log("ended")
  })
  // streamReply.on("close", () => {
  //   console.log("closed")
  // })
}

function runBidirectionalStreaming(client: GreeterClient, request: HelloRequest) {
  const stream = client.sayHelloBidirectionalStreaming()
  stream.on("data", (data) => {
    console.log("Received from server:", data.message)
  })
  stream.on("end", () => {
    console.log("ended")
  })

  const streamWrite = async (i: number) => {
    stream.write({
      message: `${request.message} no.${i}`,
    })
  }
  const streamConcurent = async () => {
    for (let i = 0; i < 10; i++) {
      await streamWrite(i + 1)
      await wait(500)
    }
    stream.end()
  }

  streamConcurent()
}

function main() {
  const args = parseArgs(process.argv.slice(2), {
    string: ["target", "cmd"],
  })

  let target = "localhost:50051"
  if (args.target) {
    target = args.target
  }

  const CMDS = ["say_hello", "stream_reply", "stream_request", "bidirection_stream"]
  let cmd = "say_hello"
  if (args.cmd && CMDS.includes(args.cmd)) {
    cmd = args.cmd
  }

  const client = new GreeterClient(target, grpc.credentials.createInsecure())
  const request: HelloRequest = {
    message: "Hello from client",
  }
  if (args._.length > 0) {
    request.message = args._[0]
  }

  switch (cmd) {
    case "say_hello": {
      runSayHello(client, request)
      break
    }
    case "stream_reply": {
      runStreamReply(client, request)
      break
    }
    case "stream_request": {
      runStreamRequest(client, request)
      break
    }
    case "bidirection_stream": {
      runBidirectionalStreaming(client, request)
      break
    }
  }
}

main()
