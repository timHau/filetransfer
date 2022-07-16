package handler

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Download struct {
	Err      string
	FileName string
}

func ParseFileName(name string) (int64, string, string) {
	parts := strings.Split(name, "____")
	if len(parts) != 3 {
		return 0, "", ""
	}

	t, _ := strconv.ParseInt(parts[0], 10, 64)
	return t, parts[1], parts[2]
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing file name", http.StatusBadRequest)
		return
	}

	fp := path.Join("templates", "download.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, _, name := ParseFileName(fileName)
	date := time.Unix(t, 0)
	lastMonth := time.Now().Add(-30 * 24 * time.Hour)
	if date.Before(lastMonth) {
		if err := tmpl.Execute(w, Download{Err: "File is too old"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	filePath := path.Join("assets", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		if err := tmpl.Execute(w, Download{Err: "File not found"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	http.ServeContent(w, r, name, date, file)
}
