grpc-go-stat-service
===

grpc-go-stat-service is runtime of golang information service.

inspired by [golang-stats-api-handler](golang-stats-api-handler)

# Installation

```console
$ go get github.com/f110/grpc-go-stat-service
```

# Usage

Server:

```go
import (
	"github.com/f110/grpc-go-stat-service"
	"github.com/golang/protobuf/ptypes/any"
)

func NewServer() {
	s := statservice.New(func () *any.Any {return nil})
	statservice.RegisterStatServer(grpc.Server, statservice.StatServer(s))
}
```

Client:

```go
import (
	"github.com/f110/grpc-go-stat-service"
	"google.golang.org/grpc"
)

func NewClient() {
    conn, _ := grpc.Dial(ServerAddr)
    statClient := statservice.NewStatClient(conn)
    stat := statClient.Get(context.Background(), &statservice.GetRequest{})
}
```

[golang-stats-api-handler]:https://github.com/fukata/golang-stats-api-handler