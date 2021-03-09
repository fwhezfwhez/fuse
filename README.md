## Fuse
an api fuse realized in golang.

## Start
`go get github.com/fwhezfwhez/fuse`


## Usage
### 1. fuse rpc calling

- Know that keep aware of what kind of errors should count fs.Fail() times.
```go
// 20s allows 5 fail times, on arriving 5 times will fall into fused state for 10s.
fs := NewFuse(5, 10, 20, 1)
var rpcName = "/get-user-info/"

if !fs.FuseOK(rpcName) {
    return fmt.Errorf("'%s' is fused for breaking max times", rpcName)
}

if e:=rpcCall(); e!=nil {
    fs.Fail(rpcName)
    return
}

```

### 2. fuse for gateway

**http-gin**

```
import "github.com/fwhezfwhez/fuse/middleware"
...
r:=gin.Default()
r.Use(middleware.GinHTTPFuse)
```
