package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"path"
	"time"
)

func sendMail(to string, fileName string) error {
	from := os.Getenv("MAIL_ADDRESS")
	password := os.Getenv("MAIL_PASSWORD")
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")

	fmt.Println(from, password, smtpHost, smtpPort)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	fp := path.Join("templates", "mail.html")
	t, err := template.ParseFiles(fp)
	if err != nil {
		return err
	}

	var msg bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg.Write([]byte(fmt.Sprintf("Subject: Filetransfer \n%s\n\n", mimeHeaders)))

	t.Execute(&msg, struct {
		Url string
	}{
		Url: fileName,
	})

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg.Bytes())
	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{"tim.hau@hotmail.de"}, []byte("Some one uses your filetransfer: "+to))
	if err != nil {
		return err
	}

	return nil
}

func hashedFileName(name string) string {
	hash := sha256.Sum256([]byte(name))
	t := time.Now().Unix()
	return fmt.Sprintf("%v____%x____%s", t, hash, name)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	to := r.FormValue("to")
	if to == "" || !validEmail(to) {
		http.Error(w, "Missing to", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	name := hashedFileName(handler.Filename)

	go func() {
		fileUrl := os.Getenv("SERVER_URL") + "/download?file=" + name
		err := sendMail(to, fileUrl)
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
