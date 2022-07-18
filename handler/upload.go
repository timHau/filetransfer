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

	recipient := r.FormValue("recipient")
	if recipient == "" || !utils.ValidEmail(recipient) {
		http.Error(w, "missing or invalid recipient", http.StatusBadRequest)
		return
	}

	sender := r.FormValue("sender")
	if sender == "" || !utils.ValidEmail(sender) {
		http.Error(w, "missing or invalid sender", http.StatusBadRequest)
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
		err := utils.SendMail(sender, recipient, name)
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
