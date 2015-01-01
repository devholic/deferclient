# Defer Panic client
[![GoDoc](https://godoc.org/github.com/deferpanic/deferclient?status.svg)](https://godoc.org/github.com/deferpanic/deferclient)

[![wercker status](https://app.wercker.com/status/b7a471949687969984843f7c5e5988a2/s "wercker status")](https://app.wercker.com/project/bykey/b7a471949687969984843f7c5e5988a2)

Defer Panic Client Lib.

### Installation
``go get github.com/deferpanic/deferclient``


Get an API KEY:
```
 curl https://api.deferpanic.com/v1/users/create \
        -X POST \
        -d "email=test@test.com" \
        -d "password=password"
```

### Examples

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

### Documentation

See https://godoc.org/github.com/deferpanic/deferclient for documentation.

Defer Panic Client
