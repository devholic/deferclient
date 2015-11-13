package deferstats

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

type TestContextWrapper struct {
	Ctx     context.Context
	Handler ContextHandler
}

func (wrapper *TestContextWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrapper.Handler.ServeHTTPContext(wrapper.Ctx, w, r)
}

type TestContextJSON struct {
	Title string
}

func TestHTTPContextPost(t *testing.T) {
	dps := NewClient("token")

	mux := http.NewServeMux()
	mux.Handle("/", &TestContextWrapper{Ctx: context.Background(),
		Handler: dps.HTTPContextHandlerFunc(ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}

			var tj TestContextJSON
			err = json.Unmarshal(body, &tj)
			if err != nil {
				panic(err)
			}
			if tj.Title != "sample title in json" {
				t.Error("not parsing the POST body correctly")
			}
		}))})

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	var jsonStr = []byte(`{"Title":"sample title in json"}`)
	req, err := http.NewRequest("POST", lurl, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Custom-Header", "some header")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func TestHTTPContextPostHandler(t *testing.T) {
	dps := NewClient("token")

	mux := http.NewServeMux()
	post := ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var tj TestContextJSON
		err = json.Unmarshal(body, &tj)
		if err != nil {
			panic(err)
		}
		if tj.Title != "sample title in json" {
			t.Error("not parsing the POST body correctly")
		}
	})
	mux.Handle("/", &TestContextWrapper{Ctx: context.Background(), Handler: dps.HTTPContextHandlerFunc(post)})

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	var jsonStr = []byte(`{"Title":"sample title in json"}`)
	req, err := http.NewRequest("POST", lurl, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Custom-Header", "some header")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func TestHTTPContextHeader(t *testing.T) {
	dps := NewClient("token")

	mux := http.NewServeMux()
	mux.Handle("/", &TestContextWrapper{Ctx: context.Background(),
		Handler: dps.HTTPContextHandlerFunc(ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Custom-Header") != "some header" {
				t.Error("headers not being passed through")
			}

			if r.Header.Get("Content-Type") != "application/json" {
				t.Error("headers not being passed through")
			}
		}))})

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	var jsonStr = []byte(`{"Title":"sample title in json"}`)
	req, err := http.NewRequest("POST", lurl, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Custom-Header", "some header")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func TestHTTPContextHeaderHandler(t *testing.T) {
	dps := NewClient("token")

	mux := http.NewServeMux()
	post := ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "some header" {
			t.Error("headers not being passed through")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("headers not being passed through")
		}
	})
	mux.Handle("/", &TestContextWrapper{Ctx: context.Background(), Handler: dps.HTTPContextHandlerFunc(post)})

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	var jsonStr = []byte(`{"Title":"sample title in json"}`)
	req, err := http.NewRequest("POST", lurl, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Custom-Header", "some header")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func TestContextSOA(t *testing.T) {
	curlist.Reset()

	dps := NewClient("token")

	mux := http.NewServeMux()
	mux.Handle("/", &TestContextWrapper{Ctx: context.Background(),
		Handler: dps.HTTPContextHandlerFunc(ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			r.ParseForm()

			dpsi := r.Header.Get("X-dpparentspanid")
			okey := r.FormValue("other_key")

			if dpsi != "8103318854963911860" {
				t.Error("span not accessible")
			}
			if okey != "2" {
				t.Error("other_key not accessible")
			}
		}))})

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	data := url.Values{}
	data.Set("other_key", "2")

	client := &http.Client{}
	r, err := http.NewRequest("POST", lurl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Add("X-dpparentspanid", "8103318854963911860")

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if len(curlist.list) == 0 {
		t.Error("should have a http in the list")
	}

	if curlist.list[0].ParentSpanId != 8103318854963911860 {
		t.Error("not tracking our parent_span_id")
	}
}

func TestContextSOAHandler(t *testing.T) {
	curlist.Reset()

	dps := NewClient("token")

	mux := http.NewServeMux()
	post := ContextHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		dpsi := r.Header.Get("X-dpparentspanid")
		okey := r.FormValue("other_key")

		if dpsi != "8103318854963911860" {
			t.Error("span not accessible")
		}
		if okey != "2" {
			t.Error("other_key not accessible")
		}
	})

	mux.Handle("/", &TestContextWrapper{Ctx: context.Background(), Handler: dps.HTTPContextHandlerFunc(post)})

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"

	go http.Serve(l, mux)

	lurl := "http://" + l.Addr().String() + "/"

	data := url.Values{}
	data.Set("other_key", "2")

	client := &http.Client{}
	r, err := http.NewRequest("POST", lurl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Add("X-dpparentspanid", "8103318854963911860")

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if len(curlist.list) == 0 {
		t.Error("should have a http in the list")
	}

	if curlist.list[0].ParentSpanId != 8103318854963911860 {
		t.Error("not tracking our parent_span_id")
	}
}
