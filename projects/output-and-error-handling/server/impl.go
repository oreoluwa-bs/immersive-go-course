package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		randNum := rand.Intn(10)

		switch randNum {
		case 1, 2:
			w.Write([]byte("Sunny"))
		case 3, 4:
			w.Write([]byte("Rainy"))
		case 5, 6, 7, 8:
			// Retry
			secsToRetry := rand.Intn(15) + 1

			if secsToRetry > 5 {

				w.WriteHeader(503)
				w.Write([]byte("Connection timedout. Try again later"))
				return
			}

			w.WriteHeader(429)
			w.Header().Add("Retry-After", strconv.Itoa(secsToRetry))
			resp := fmt.Sprintf("Connection failed. Retrying  in %ds", secsToRetry)

			w.Write([]byte(resp))

			time.Sleep(time.Duration(secsToRetry) * time.Second)

			randNum := rand.Intn(4)

			switch randNum {
			case 3:
				w.Write([]byte("Sunny"))

			default:
				w.Write([]byte("Rainy"))
			}

		default:
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
		}
	})

	fmt.Fprintln(os.Stderr, "Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen: %v", err)
		os.Exit(1)
	}

}
