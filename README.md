### Example Server
```go
package main

import (
    "log"

    "github.com/TheMrViper/groto"
)

type Request struct {
    Field1 int
    Field2 int
}
type Response struct {
    Result int
}

func (*Request) Name() string {
    return "testRequest"
}
func (*Request) Create() groto.ServerRequest {
    return &Request{}
}
func (req *Request) Handler() groto.ServerResponse {
    res := &Response{}
    res.Result = req.Field1 + req.Field2
    return res
}

func main() {
    serv := groto.NewServer(":8080")
    
    serv.Handle(&Request{})
    
    err := serv.Listen()
    if err != nil {
        log.Fatalln(err)
    }
}
```

### Example Client
```go
package main 

import (
    "log"
    
    "github.com/TheMrViper/groto"
)

type Request struct {
    Field1 int
    Field2 int
}
type Response struct {
    Result int
}

func (*Request) Name() string {
    return "testRequest"
}

func main() {
    cli := groto.NewClient("127.0.0.1:8080", 1)
    
    req := &Request{}
    req.Field1 = 1
    req.Field2 = 2
    
    res := &Response{}
    
    err := cli.Send(req, res)
    if err != nil {
        log.Fatalln(err)
    }
    log.Println("Result: ", res.Result)
}
```

### Recommended Server hierarchy 
```
ServiceDirectory 
    req
        ServerRequestName.go
    res
        ServerResponseName.go
    main.go
```