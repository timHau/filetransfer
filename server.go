package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/timHau/filetransfer/handler"
	"github.com/timHau/filetransfer/jobs"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jobs.DeleteJob()
	m := jobs.MultipleJob()
	go m.Run()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/upload", handler.HandleUpload)
	http.HandleFunc("/multi", func(w http.ResponseWriter, r *http.Request) { handler.HandleMulti(w, r, m) })
	http.HandleFunc("/download", handler.HandleDownload)
	http.HandleFunc("/", handler.HandleSite)

	fmt.Printf("Starting server on port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
