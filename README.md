# protoc-gen-go-asynq

generate protobuf for asynq payload.

## usage

### 1. install protoc-gen-go-asynq

```bash
# first install protobuf, then install
go get -d google.golang.org/protobuf/cmd/protoc-gen-go
go get -d google.golang.org/grpc/cmd/protoc-gen-go-grpc

# install protoc-gen-go-asynq
go get -d github.com/jmploop/protoc-gen-go-asynq@v1.0.2
```

### 2. copy proto/asynq.proto to your project

### 3. define task proto

```protobuf
syntax = "proto3";

package example;

// import copy asynq.proto
import "asynq.proto";
import "google/protobuf/empty.proto";

option go_package = ".;example";
option java_multiple_files = true;
option java_package = "example";

service User {
    rpc CreateUser(CreateUserPayload) returns (google.protobuf.Empty) {
        option (asynq.task) = {
            typename: "user:create" // define asynq.task typename, it is unique.
        };
    };
    rpc UpdateUser(UpdateUserPayload) returns (google.protobuf.Empty) {
        option (asynq.task) = {
            typename: "user:update"
        };
    };
}

message CreateUserPayload {
    string name = 1;
}

message UpdateUserPayload {
    string name = 1;
}
```

### 4. generate

```bash
protoc --proto_path=. \
       --proto_path=./{copy_dir} \
       --go_out=paths=source_relative:. \
       --go-asynq_out=paths=source_relative:. \
       example.proto
```

### 5. implement

Don't care about encode/decode, you just implement your define method.
