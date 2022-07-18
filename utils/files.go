package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func MergeMultiFiles(name string) error {
	dirPath := path.Join("./assets", name)
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	numOfMulti, _ := strconv.ParseInt(os.Getenv("NUM_OF_MULTI"), 10, 64)
	fileNames := make([]string, numOfMulti)

	for _, file := range files {
		num, _, err := ParseMultiFile(file.Name())
		if err != nil {
			log.Println("Error while parsing", file.Name(), err)
			continue
		}

		fileNames[num] = file.Name()
	}

	out, err := os.Create(path.Join("./assets", HashedFileName(name)))
	if err != nil {
		return err
	}

	for _, fileName := range fileNames {
		filePath := path.Join("./assets", fileName)
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		io.Copy(out, file)
	}

	if err = os.RemoveAll(dirPath); err != nil {
		return err
	}

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
