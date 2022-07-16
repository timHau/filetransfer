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
	fmt.Println("[Handling files]")

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

func Init() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Seconds().Do(handleFile)
	s.StartAsync()
}
