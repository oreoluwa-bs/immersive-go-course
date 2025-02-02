package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Image struct {
	Title   string
	AltText string
	URL     string
}

func main() {
	// images := []Image{
	// 	{
	// 		Title:   "Sunset",
	// 		AltText: "Clouds at sunset",
	// 		URL:     "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
	// 	},
	// 	{
	// 		Title:   "Mountain",
	// 		AltText: "A mountain at sunset",
	// 		URL:     "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
	// 	},
	// }

	conn := connectToDB()

	defer conn.Close(context.Background())

	mux := http.NewServeMux()

	mux.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		indent := 2

		indQuery := r.URL.Query().Get("indent")
		if indQuery != "" {
			if v, err := strconv.Atoi(indQuery); err == nil {
				indent = v
			}
		}

		var images []Image

		if r.Method == "GET" {
			ims, err := fetchImages(conn)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			images = ims
		}

		if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}

			var newImg Image
			err = json.Unmarshal(body, &newImg)
			if err != nil {
				http.Error(w, "Error parsing JSON", http.StatusBadRequest)
				return
			}

			ims, err := addImage(conn, newImg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			images = ims
		}

		data, err := json.MarshalIndent(images, "", strings.Repeat(" ", indent))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	http.ListenAndServe(":8000", mux)
}

func connectToDB() *pgx.Conn {

	DB_URL := os.Getenv("DATABASE_URL")

	if DB_URL == "" {
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// defer conn.Close(context.Background())

	return conn
}

func fetchImages(conn *pgx.Conn) ([]Image, error) {

	rows, err := conn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")

	if err != nil {
		return nil, err
	}

	images := make([]Image, 0)

	for rows.Next() {
		var image Image

		err := rows.Scan(&image.Title, &image.URL, &image.AltText)

		if err != nil {
			return nil, err
		}

		images = append(images, image)

	}

	return images, nil
}

func addImage(conn *pgx.Conn, i Image) ([]Image, error) {

	_, err := conn.Exec(context.Background(), "INSERT INTO public.images(title,url,alt_text) VALUES($1, $2, $3)", i.Title, i.URL, i.AltText)

	if err != nil {
		return nil, err
	}

	images := make([]Image, 0)

	images = append(images, i)

	return images, nil
}
