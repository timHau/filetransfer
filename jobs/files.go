package jobs

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/timHau/filetransfer/handler"
)

func handleFile() {
	files, _ := ioutil.ReadDir("./assets")
	for _, f := range files {
		if f.Name()[0:1] == "." {
			continue
		}

		t, _, _ := handler.ParseFileName(f.Name())
		if time.Unix(t, 0).Add(time.Hour * 24).Before(time.Now()) {
			fmt.Printf("[Deleting file] %s\n", f.Name())
			os.Remove("./assets/" + f.Name())
		}
	}
}

func DeleteJob() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Hour().Do(handleFile)
	s.StartAsync()
}
