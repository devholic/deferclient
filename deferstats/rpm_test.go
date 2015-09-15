package deferstats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

func TestRPM(t *testing.T) {

	dps := NewClient("token")

	mux := http.NewServeMux()
	mux.HandleFunc("/200", dps.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var tj TestJSON
		err = json.Unmarshal(body, &tj)
		if tj.Title != "sample title in json" {
			t.Error("not parsing the POST body correctly")
		}

		fmt.Fprintf(w, "200")
	}))

	mux.HandleFunc("/500", dps.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("500!!!")
		fmt.Fprintf(w, "200")
	}))

	// set listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error("http not listening")
	}

	dps.statsUrl = "http://" + l.Addr().String() + "/"
	dps.BaseClient.NoPost = true

	go http.Serve(l, mux)

	url200 := "http://" + l.Addr().String() + "/200"
	url500 := "http://" + l.Addr().String() + "/500"

	rpmz := rpms.List()
	if rpmz.StatusOk != 0 {
		t.Errorf("StatusOk is not 0 %v", rpmz.StatusOk)
	}

	var jsonStr = []byte(`{"Title":"sample title in json"}`)
	for i := 0; i < 3; i++ {

		req, err := http.NewRequest("POST", url200, bytes.NewBuffer(jsonStr))
		req.Header.Set("X-Custom-Header", "some header")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}

	req, err := http.NewRequest("POST", url500, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "some header")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	rpmz = rpms.List()
	if rpmz.StatusOk != 3 {
		t.Errorf("not inc'ing StatusOk %v", rpmz.StatusOk)
	}

	if rpmz.StatusInternalServerError != 1 {
		t.Errorf("not inc'ing StatusInternalServerError %v", rpmz.StatusInternalServerError)
	}

}

func TestClearRPM(t *testing.T) {
	rpms.ResetRPM()

	rpms.Inc(200)

	if rpms.rpm.StatusOk != 1 {
		t.Errorf("not inc'ing StatusOk %v", rpms.rpm.StatusOk)
	}

	rpms.ResetRPM()

	if rpms.rpm.StatusOk != 0 {
		t.Errorf("not clearing StatusOk %v", rpms.rpm.StatusOk)
	}

}
