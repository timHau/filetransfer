package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/timHau/filetransfer/jobs"
	"github.com/timHau/filetransfer/utils"
)

func HandleMulti(w http.ResponseWriter, r *http.Request, m *jobs.Multiple) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024*1024)

	to := r.FormValue("to")
	if to == "" || !utils.ValidEmail(to) {
		http.Error(w, "Missing to", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error while getting file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, fileName, err := utils.ParseMultiFile(handler.Filename)
	if err != nil {
		log.Println("Error while parse", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dirPath := path.Join("./assets", fileName)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, 0777)
	}

	f, err := os.OpenFile(dirPath+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)

	m.Receiver <- utils.FileUploadMessage{
		Name:    fileName,
		EmailTo: to,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
