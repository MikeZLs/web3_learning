
#### errors.go 源码
```go
package builtin

type error interface {
	Error() string
}
```

```go
package errors

func New(text string) error {
    return &errorString{text}
}

type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
```

#### 自定义一个错误类型
```go
type MyError struct {
    Code int
    Msg  string
}

func (me MyError) Error() string {
    return fmt.Sprintf("code=%d, msg=%s", me.Code, me.Msg)
}

func doingError() error {
    return MyError{Code: 909, Msg: "未知错误！"}
}
```
