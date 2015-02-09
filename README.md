# Defer Panic client
[![GoDoc](https://godoc.org/github.com/deferpanic/deferclient?status.svg)](https://godoc.org/github.com/deferpanic/deferclient)

[![wercker status](https://app.wercker.com/status/b7a471949687969984843f7c5e5988a2/s "wercker status")](https://app.wercker.com/project/bykey/b7a471949687969984843f7c5e5988a2)

Defer Panic Client Lib.

 *  **Error Handling** - DeferClient can auto-wrap your errors to shoot
    up to deferpanic or you can choose to explicitly log the ones you
    care about.

 *  **Panic Handling** - Let deferclient catch and log any panic you get
    to your own dashboard.

 *  **HTTP latency** - DeferClient can log the latencies of all your hard
    hit http requests.

 *  **Database latency** - Get notified of slow database queries in your
    go app.

 *  **Metrics** - See goroutines, memory usage, gc and more automatically
    in your own dashboard.

 *  **Custom K/V** - Got something we don't support? You can log your own k/v metrics just as easily.


Translations:

* [简体中文](README_zh_cn.md)
* [Русский](README_ru_RU.md)

### Installation
``go get github.com/deferpanic/deferclient``

**api key**

Get an API KEY via your shell or signup manually [here](https://deferpanic.com/signup):
```
 curl https://api.deferpanic.com/v1/users/create \
        -X POST \
        -d "email=test@test.com" \
        -d "password=password"
```

### HTTP Examples

Here we have 4 examples:
* log a fast request
* log a slow request
* log an error
* log a panic

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

The client works perfectly fine in non-HTTP applications:

### Non-HTTP Errors/Panics - automatic errors
Here we log both an error and a panic. If you want us to catch your
errors when you instantiate them use this method. We'll create a new
deferpanic error that is returned and the error will be shipped to
deferpanic.

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

### Errors - explicitly log
If you want to explicitly log your errors use this method. deferlog.Wrap
will log the bactrace and the error and ship it up to deferpanic
immediately.

```
package main

import (
        "errors"
        "fmt"
        "github.com/deferpanic/deferclient/deferlog"
        "time"
)

func errorTest() {
        err := errors.New("danger will robinson!")
        if err != nil {
                deferlog.Wrap(err)
                fmt.Println(err)
        }
}

func main() {

        deferlog.Token = "v00L0K6CdKjE4QwX5DL1iiODxovAHUfo"

        errorTest()

        time.Sleep(5 * time.Second)
}
```

### Database Latency

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

### Generic K/V
If you wish to log other k/v metrics this implements a very basic
counter over time.

```
package main

import (
        "github.com/deferpanic/deferclient/deferkv"
        "time"
)

func main() {
        deferkv.Token = "v00L0K6CdKjE4QwX5DL1iiODxovAHUfo"

        deferkv.Report("some_key", 10)

        deferkv.Report("some_other_key", 30)

        time.Sleep(5 * time.Second)
}
```

### Documentation

See https://godoc.org/github.com/deferpanic/deferclient for documentation.

Defer Panic Client
