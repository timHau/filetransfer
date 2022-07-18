package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/timHau/filetransfer/utils"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set max file size to 10GB
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024*1024)

	to := r.FormValue("to")
	if to == "" || !utils.ValidEmail(to) {
		http.Error(w, "Missing to", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	name := utils.HashedFileName(handler.Filename)

	go func() {
		fileUrl := os.Getenv("SERVER_URL") + "/download?file=" + name
		err := utils.SendMail(to, fileUrl)
		if err != nil {
			fmt.Println(err)
		}
	}()

	f, err := os.OpenFile("./assets/"+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
