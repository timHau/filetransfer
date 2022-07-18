package jobs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/timHau/filetransfer/utils"
)

func handleFile() {
	// storage time in hour
	storageTime, err := strconv.ParseInt(os.Getenv("STORAGE_TIME"), 10, 64)
	if err != nil {
		log.Fatal("Error while parsing STORAGE_TIME")
	}

	files, _ := ioutil.ReadDir("./assets")
	for _, f := range files {
		if f.Name()[0:1] == "." {
			continue
		}

		if f.IsDir() {
			// delete directory if it is not used
			// just to be sure to clean up if something went wrong while uploading
			if f.Mode().IsRegular() {
				if !f.ModTime().Add(time.Hour * time.Duration(storageTime)).Before(time.Now()) {
					if err := os.RemoveAll(f.Name()); err != nil {
						log.Println("Error while removing", f.Name(), err)
					}
				}
			}
		} else {
			t, _, _ := utils.ParseFileName(f.Name())
			if time.Unix(t, 0).Add(time.Hour * time.Duration(storageTime)).Before(time.Now()) {
				fmt.Printf("[Deleting file] %s\n", f.Name())
				os.Remove("./assets/" + f.Name())
			}
		}

	}
}

func DeleteJob() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Hour().Do(handleFile)
	s.StartAsync()
}
