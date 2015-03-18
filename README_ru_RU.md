# Defer Panic клиент
[![GoDoc](https://godoc.org/github.com/deferpanic/deferclient?status.svg)](https://godoc.org/github.com/deferpanic/deferclient)

[![wercker status](https://app.wercker.com/status/b7a471949687969984843f7c5e5988a2/s "wercker status")](https://app.wercker.com/project/bykey/b7a471949687969984843f7c5e5988a2)

Клиент Defer Panic.

### установка
``go get github.com/deferpanic/deferclient``


Получение ключа API:
```
 curl https://api.deferpanic.com/v1/users/create \
        -X POST \
        -d "email=test@test.com" \
        -d "password=password"
```

### Примеры с net/http

Примеры работы с Defer Panic:
* Лог быстрых запросов
* Лог медленных запросов
* Журналирование ошибки
* Журналирование паники

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

Также клиент может использоваться и в других приложениях:

### Автоматическое логирование паник/ошибок
В данном примере мы отправляем в лог и ошибку и панику. 

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

### Медленные запросы к БД

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

### Микросервисы/SOA

Используете микросервисы? Теперь вы можете отследить ваши субзапросы 
через всё приложение до корневого запроса!

Пример

Публичный сервис
```
package main

import (
    "fmt"
    "github.com/deferpanic/deferclient/deferstats"
    "io/ioutil"
    "net/http"
    "net/url"
)

func handler(w http.ResponseWriter, r *http.Request) {

    // just pass your spanId w/each request
    resp, err := http.PostForm("http://127.0.0.1:7070/internal",
        url.Values{"defer_parent_span_id": {deferstats.GetSpanIdString(w)}})
    if err != nil {
        fmt.Println(err)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    fmt.Fprintf(w, string(body))
}

func main() {
    deferstats.Token = "v00L0K6CdKjE4QwX5DL1iiODxovAHUfo"

    go deferstats.CaptureStats()

    http.HandleFunc("/", deferstats.HTTPHandler(handler))
    http.ListenAndServe(":9090", nil)
}
```

Внутренний API
```
package main

import (
        "encoding/json"
        "github.com/deferpanic/deferclient/deferstats"
        "net/http"
        "time"
)

type blah struct {
        Stuff string
}

func handler(w http.ResponseWriter, r *http.Request) {
        time.Sleep(250 * time.Millisecond)

        stuff := blah{
                Stuff: "some reply",
        }

        js, err := json.Marshal(stuff)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Write(js)
}

func main() {
        deferstats.Token = "v00L0K6CdKjE4QwX5DL1iiODxovAHUfo"

        go deferstats.CaptureStats()

        http.HandleFunc("/internal", deferstats.HTTPHandler(handler))
        http.ListenAndServe(":7070", nil)
}
```

При обёртке стандартного net/http handler'а в deferstats.HTTPHandler
публичный сервис автоматически привязывается к внутреннему.

### Документация

См. [документацию на GoDoc](https://godoc.org/github.com/deferpanic/deferclient).

Клиент Defer Panic
