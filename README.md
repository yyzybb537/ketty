# ketty

### INSTALL

```
$ go get github.com/yyzybb537/{ketty,kettygen}
```

### Generate Code

```
$ kettygen xxx.proto
```

### Example

#### Proto
```
syntax = "proto3";

package test_pb;

message Req {
    int64 val = 1;
}

message Rsp {
    int64 val = 1;
}

service EchoService {
    rpc Echo(Req) returns(Rsp) {}
}
```

#### Server
```
package main

import (
        echo "github.com/yyzybb537/ketty/demo/pb"
        "github.com/yyzybb537/ketty"
        "golang.org/x/net/context"
)

type EchoServer struct {
}

func (this *EchoServer) Echo(ctx context.Context, req *echo.Req) (*echo.Rsp, error) {
        return &echo.Rsp{Val:req.Val}, nil
}

func main() {
        server, err := ketty.Listen("grpc://127.0.0.1:8090", "")
        if err != nil {
                ketty.GetLog().Errorf("Listen error:%s", err.Error())
        }

        server.RegisterMethod(echo.EchoServiceHandle, &EchoServer{})

        err = server.Serve()
        if err != nil {
                ketty.GetLog().Errorf("Serve error:%s", err.Error())
        }

        q := make(chan int)
        <-q
}
```

#### Client
```
package main

import (
        echo "github.com/yyzybb537/ketty/demo/pb"
        "github.com/yyzybb537/ketty"
        "golang.org/x/net/context"
)

func main() {
        client, err := ketty.Dial("grpc://127.0.0.1:8090", "")
        if err != nil {
                ketty.GetLog().Errorf("Dial error:%s", err.Error())
                return
        }
        defer client.Close()

        req := &echo.Req{ Val : 123 }
        stub := echo.NewKettyEchoServiceClient(client)
        rsp, err := stub.Echo(context.Background(), req)
        if err != nil {
                ketty.GetLog().Errorf("Invoke error:%+v", err)
                return
        }

        ketty.GetLog().Infof("Rsp: %d", rsp.Val)
}
```

### URL


```
// @sUrl:   protocol://ip[:port][,ip[:port]]/path                                  
// E.g:                                                                            
//    http://127.0.0.1:8030/path                                                   
//    https://127.0.0.1:8030                                                     
//    grpc://127.0.0.1:8030,127.0.0.1:8031                                         
//                                                                                 
// @sBalanceUrl:  driver://ip[:port][,ip[:port]]/path                            
//    etcd://127.0.0.1:2379/path   
```
