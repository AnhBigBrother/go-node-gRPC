import parseArgs from "minimist"
import { GreeterClient, HelloRequest, HelloReply } from "./proto/helloworld"
import * as grpc from "@grpc/grpc-js"
import { readFileSync } from "node:fs"

async function wait(time: number) {
  return new Promise((resolve) => setTimeout(resolve, time))
}

function runSayHello(client: GreeterClient, request: HelloRequest) {
  const metaData = new grpc.Metadata()
  metaData.add("authorization", "supersecret")
  metaData.add("description", "this is a unary request")

  client.sayHello(
    request,
    metaData,
    function (err: grpc.ServiceError | null, response: HelloReply) {
      console.log("Received from server:", response?.message)
    }
  )
}

function runStreamRequest(client: GreeterClient, request: HelloRequest) {
  const metaData = new grpc.Metadata()
  metaData.add("authorization", "supersecret")
  metaData.add("description", "this is a stream request")

  const streamRequest = client.sayHelloStreamRequest(
    metaData,
    (err: grpc.ServiceError | null, response: HelloReply) => {
      if (err) {
        console.log(err)
        return
      }
      console.log("Received from server:", response.message)
    }
  )

  const writeAsync = async (num: number) => {
    streamRequest.write({
      message: `${request.message} no.${num}`,
    })
  }
  const streamAsync = async () => {
    for (let i = 0; i < 10; i++) {
      await writeAsync(i + 1)
      await wait(500)
    }
    streamRequest.end()
  }

  streamAsync()
}

function runStreamReply(client: GreeterClient, request: HelloRequest) {
  const metaData = new grpc.Metadata()
  metaData.add("authorization", "supersecret")
  metaData.add("description", "this is a stream request")

  const streamReply = client.sayHelloStreamReply(request, metaData)
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
  const metaData = new grpc.Metadata()
  metaData.add("authorization", "supersecret")
  metaData.add("description", "this is a bidirection stream request")

  const stream = client.sayHelloBidirectionalStreaming(metaData)
  stream.on("data", (data) => {
    console.log("Received from server:", data.message)
  })
  stream.on("end", () => {
    console.log("ended")
  })

  const writeAsync = async (i: number) => {
    stream.write({
      message: `${request.message} no.${i}`,
    })
  }
  const streamAsync = async () => {
    for (let i = 0; i < 10; i++) {
      await writeAsync(i + 1)
      await wait(500)
    }
    stream.end()
  }

  streamAsync()
}

function loadCreadential(): grpc.ChannelCredentials {
  const rootCert = readFileSync("../cert/ca-cert.pem")
  return grpc.credentials.createSsl(rootCert)
}

function main() {
  const args = parseArgs(process.argv.slice(2), {
    string: ["target", "cmd"],
  })

  let target = "localhost:8080"
  if (args.target) {
    target = args.target
  }

  const CMDS = ["say_hello", "stream_reply", "stream_request", "bidirection_stream"]
  let cmd = "say_hello"
  if (args.cmd && CMDS.includes(args.cmd)) {
    cmd = args.cmd
  }

  const client = new GreeterClient(target, loadCreadential())
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
