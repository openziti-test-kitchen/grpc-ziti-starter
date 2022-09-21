OpenZiti gRPC Project Template
==============================

Use this project to quickly start your next gRPC project that uses open-source, secure, zero-trust
[OpenZiti](https://openziti.io) Network.


## Start

Create your project's repo using this one as a 
[template](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template).

## Try it

* Get yourself an OpenZiti network:
  * follow [quickstart](https://openziti.github.io/ziti/quickstarts/quickstart-overview.html) docs
  * or, use Ziti Edge Developer Sandobox([ZEDS](https://zeds.openziti.org))
* Create a Ziti [service](https://openziti.github.io/ziti/services/overview.html) to use for gRPC.
* Create, and enroll you server and client [identities](https://openziti.github.io/ziti/identities/overview.html)
* Run it! See below:
  * 'grpc-service' - name of the service
  * 'client.json' - name of client identity file
  * 'server.json' - name of server identity file

Run server:
```console
$ go run ./server -identity server.json -service grpc-service
2022/09/14 14:27:17 server listening at grpc-service
...
```
Note: the server is powered by OpenZiti network and does not have any inbound ports open

Run client:
```console
$ go run ./client -identity client.json -service grpc-service -what-is foo
2022/09/14 14:29:32 Answer: I don't know what foo is :(

$ go run ./client -identity client.json -service grpc-service -what-is ziti
2022/09/14 14:29:56 Answer: ziti is a type of pasta
```

## Next Steps

* design and implement your gRPC API  
  * modify protocol/starter.proto to fit your needs
  * generate Golang code for the protocol: `$ protoc --go_out=plugins=grpc:. ./protocol/starter.proto`
  * make changes in `server/` package to implement your API
  * make changes in `client/` package to use your new API
* deploy your server on your production OpenZiti network
* profit!!


## How it's done
In this project we use `google.golang.com/grpc`.

### Server
We start gRPC server with Ziti listener.

This is all is needed to zitify gRPC server (error handling is stripped for brevity):
```go
// bootstrap Ziti
ztx := ziti.NewContextWithConfig(cfg)
_ = ztx.Authenticate()
lis, _ := ztx.Listen(*service)

// standard gRCP init
s := grpc.NewServer()
protocol.RegisterAnswerServiceServer(s, &server{})

// serve using Ziti server connection
_ = s.Serve(lis)
```

### Client
We create gRPC client with Ziti connection, like this
```golang
// bootstrap Ziti
ztx := ziti.NewContextWithConfig(cfg)
_ = ztx.Authenticate()

// Provide Ziti Dialer to connect to ziti service
conn, err := grpc.Dial(*service,
     grpc.WithTransportCredentials(insecure.NewCredentials()),
     grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
              return ztx.Dial(s)
     }),
)

// create client
c := protocol.NewAnswerServiceClient(conn)
```


## Have questions?

* Follow our [Blog](https://openziti.io/)
* Read [Documentation](https://openziti.github.io)
* Join [Discussion](https://openziti.discourse.group)
* Contribute to [Development](https://github.com/openziti)
* Like it? Give us a [star](https://github.com/openziti/ziti)