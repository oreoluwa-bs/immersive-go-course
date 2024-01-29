package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
)

const tmpl = `
<!DOCTYPE html>
<html>
<em>Hello, world</em>
<p>Query parameters:
<ul>
{{ range .Queries}}
<li>{{.Key}}: {{.Value}}</li>
{{ end }}
</ul>
`

// Take a rate.Limiter instance and a http.HandlerFunc and return another http.HandlerFunc that
// checks the rate limiter using `Allow()` before calling the supplied handler. If the request
// is not allowed by the limiter, a `503 Service Unavailable` Error is returned.
func rateLimit(limiter *rate.Limiter, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		}
	})
}

type Query struct {
	Key   string
	Value []string
}

func main() {
	limiter := rate.NewLimiter(100, 30)

	mux := http.NewServeMux()

	// Define your routes
	mux.HandleFunc("/200", handle200)
	mux.HandleFunc("/300", handle300)
	mux.HandleFunc("/400", handle400)
	mux.HandleFunc("/500", handle500)
	mux.HandleFunc("/authenticated", handleAuthenticated)

	mux.HandleFunc("/", rateLimit(limiter, func(w http.ResponseWriter, r *http.Request) {

		type Data struct {
			Queries []Query
		}

		qs := make([]Query, 0)

		for key, values := range r.URL.Query() {
			qs = append(qs, Query{
				Value: values,
				Key:   key,
			})
		}

		t, err := template.New("query list").Parse(tmpl)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		// fmt.Println("Hi")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-type", "text/html")
		err = t.Execute(w, Data{
			Queries: qs,
		})
		if err != nil {
			log.Panic(err)
		}
	}))

	server := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func handle200(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func handle300(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusPermanentRedirect)
	w.Write([]byte("Redirected"))
}

func handle400(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad request"))
}
func handle500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}

func handleAuthenticated(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header["Authorization"]
	if len(authHeader) < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not authorized"))
		return
	}
	spl := strings.Split(authHeader[0], "Basic")
	if len(spl) < 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not authorized"))
		return
	}

	token := strings.TrimSpace(spl[1])

	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not authorized"))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
	<html>
	<em>Hello, world</em
	</ul>`))

}
