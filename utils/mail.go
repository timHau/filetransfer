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

func SendMail(sender string, recipient string, name string) error {
	fileName := os.Getenv("SERVER_URL") + "/download?file=" + name

	from := os.Getenv("MAIL_ADDRESS")
	password := os.Getenv("MAIL_PASSWORD")
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	fp := path.Join("static", "mail.html")
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

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipient}, msg.Bytes())
	if err != nil {
		return err
	}

	var msgToMe bytes.Buffer
	msgToMe.Write([]byte(fmt.Sprintf("Subject: Filetransfer Usage %s \n%s\n\n", recipient, mimeHeaders)))
	msgToMe.Write([]byte(fmt.Sprintf("Email from: %s to: %s\n\n", sender, recipient)))
	msgToMe.Write([]byte(fileName))

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{"tim.hau@hotmail.de"}, msgToMe.Bytes())
	if err != nil {
		return err
	}

	return nil
}
