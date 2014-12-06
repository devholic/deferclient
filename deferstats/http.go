package deferstats

func deferHTTPHandler(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		f(w, r)
		endTime := time.Now()
	}
}
