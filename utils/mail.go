package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"
	"os"
	"path"
)

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func SendMail(to string, fileName string) error {
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

	var msgToMe bytes.Buffer
	msgToMe.Write([]byte(fmt.Sprintf("Subject: Filetransfer Usage %s \n%s\n\n", to, mimeHeaders)))
	msgToMe.Write([]byte(to))

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{"tim.hau@hotmail.de"}, msgToMe.Bytes())
	if err != nil {
		return err
	}

	return nil
}
