package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func MergeMultiFiles(dirPath string) error {
	// read files from dirPath
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	numOfMulti, _ := strconv.ParseInt(os.Getenv("NUM_OF_MULTI"), 10, 64)
	fileNames := make([]string, numOfMulti)

	for _, file := range files {
		parts := strings.Split(file.Name(), "____")
		if len(parts) != 2 {
			continue
		}

		num, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue
		}

		log.Println("Num", num)
		fileNames[num] = file.Name()
	}

	log.Println("fileNames", fileNames)

	return nil
}

func ParseFileName(name string) (int64, string, error) {
	parts := strings.Split(name, "____")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid file name")
	}

	t, _ := strconv.ParseInt(parts[0], 10, 64)
	return t, parts[1], nil
}

func HashedFileName(name string) string {
	t := time.Now().Unix()
	return fmt.Sprintf("%v____%s", t, name)
}

func ParseMultiFile(name string) (int64, string, error) {
	parts := strings.Split(name, "____")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid file name")
	}

	num, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", err
	}

	return num, parts[1], nil
}

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
