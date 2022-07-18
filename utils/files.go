package utils

import (
	"fmt"
	"log"
	"os"
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
