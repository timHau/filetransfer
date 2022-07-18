package jobs

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/timHau/filetransfer/utils"
)

type Multiple struct {
	mu              sync.Mutex
	transferedFiles map[string]int64
	Receiver        chan utils.FileUploadMessage
	numOfMulti      int64
}

func MultipleJob() *Multiple {
	numOfMulti, _ := strconv.ParseInt(os.Getenv("NUM_OF_MULTI"), 10, 64)
	return &Multiple{
		mu:              sync.Mutex{},
		transferedFiles: make(map[string]int64),
		Receiver:        make(chan utils.FileUploadMessage),
		numOfMulti:      numOfMulti,
	}
}

func (m *Multiple) HandleReceive(fm utils.FileUploadMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transferedFiles[fm.Name]++

	if m.transferedFiles[fm.Name] == m.numOfMulti {
		if err := utils.MergeMultiFiles(fm); err != nil {
			log.Println("Error while merging", err)
		}
		delete(m.transferedFiles, fm.Name)
	}
}

func (m *Multiple) Run() {
	for msg := range m.Receiver {
		m.HandleReceive(msg)
	}
}
