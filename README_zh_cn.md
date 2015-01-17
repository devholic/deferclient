# Defer Panic 客户
[![GoDoc](https://godoc.org/github.com/deferpanic/deferclient?status.svg)](https://godoc.org/github.com/deferpanic/deferclient)

[![wercker status](https://app.wercker.com/status/b7a471949687969984843f7c5e5988a2/s "wercker status")](https://app.wercker.com/project/bykey/b7a471949687969984843f7c5e5988a2)

Defer Panic 客户端库.

### 安装
``go get github.com/deferpanic/deferclient``


得到一个API KEY:
```
 curl https://api.deferpanic.com/v1/users/create \
        -X POST \
        -d "email=test@test.com" \
        -d "password=password"
```

### HTTP范例

在这里，我们有4个例子：
*日志快速要求
*登录缓慢的请求
*记录错误
*登录恐慌

```go
package main

import (
        "fmt"
        "github.com/deferpanic/deferclient/deferstats"
        "github.com/deferpanic/deferclient/errors"
        "net/http"
        "time"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
        err := errors.New("throwing that error")
        if err != nil {
                fmt.Println(err)
        }

        fmt.Fprintf(w, "Hi")
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
        panic("there is no need to panic")
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "this request is fast")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
        time.Sleep(3 * time.Second)
        fmt.Fprintf(w, "this request is slow")
}

func main() {
        deferstats.Token = "v00L0K6CdKjE4QwX5DL1iiODxovAHUfo"

        go deferstats.CaptureStats()

        http.HandleFunc("/fast", deferstats.HTTPHandler(fastHandler))
        http.HandleFunc("/slow", deferstats.HTTPHandler(slowHandler))
        http.HandleFunc("/panic", deferstats.HTTPHandler(panicHandler))
        http.HandleFunc("/error", deferstats.HTTPHandler(errorHandler))

        http.ListenAndServe(":3000", nil)
}
```

客户端工作完全正常的非HTTP应用程序：

### 非HTTP错误/恐慌
在这里，我们记录这两个错误和恐慌。

```
package main

import (
        "fmt"
        "github.com/deferpanic/deferclient/deferclient"
        "github.com/deferpanic/deferclient/errors"
        "time"
)

func errorTest() {
        err := errors.New("erroring out!")
        if err != nil {
                fmt.Println(err)
        }
}

func panicTest() {
        defer deferclient.Persist()
        panic("there is no need to panic")
}

func main() {
        deferclient.Token = "v00L0K6CdKjE4QwX5DL1iiODxovAHUfo"

        errorTest()
        panicTest()

        time.Sleep(time.Second * 20)
}
```

### 数据库延迟

```
package main

import (
        "database/sql"
        "github.com/deferpanic/deferclient/deferstats"
        _ "github.com/lib/pq"
        "log"
        "time"
)

func main() {
        deferstats.Token = "v00L0K6CdKjE4QwX5DL1iiODxovAHUfo"

        _db, err := sql.Open("postgres", "dbname=dptest sslmode=disable")
        db := deferstats.NewDB(_db)

        go deferstats.CaptureStats()

        var id int
        var sleep string
        err = db.QueryRow("select 1 as num, pg_sleep(0.25)").Scan(&id, &sleep)
        if err != nil {
                log.Println("oh no!")
        }

        err = db.QueryRow("select 1 as num, pg_sleep(2)").Scan(&id, &sleep)
        if err != nil {
                log.Println("oh no!")
        }

        time.Sleep(3 * time.Second)
}
```

### 文件

见https://godoc.org/github.com/deferpanic/deferclient 的文档。

See https://godoc.org/github.com/deferpanic/deferclient for documentation.

Defer Panic 客户
