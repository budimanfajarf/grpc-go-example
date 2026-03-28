# gRPC Hello World

Follow these setup to run the [quick start][] example:

1.  Run the server:

    ```console
    go run greeter_server/main.go
    ```

2.  Run the client:

    ```console
    go run greeter_client/main.go
    Greeting: Hello world
    ```

3.  Regenerate gRPC code from the .proto file:

    ```console
    protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld/helloworld.proto
    ```

For more details (including instructions for making a small change to the
example code) or if you're having trouble running this example, see [Quick
Start][].

[quick start]: https://grpc.io/docs/languages/go/quickstart
