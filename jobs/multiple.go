package jobs

import (
	"os"
	"strconv"
	"sync"

	"github.com/timHau/filetransfer/utils"
)

type Multiple struct {
	mu              sync.Mutex
	transferedFiles map[string]int64
	Receiver        chan string
	numOfMulti      int64
}

func MultipleJob() *Multiple {
	numOfMulti, _ := strconv.ParseInt(os.Getenv("NUM_OF_MULTI"), 10, 64)
	return &Multiple{
		mu:              sync.Mutex{},
		transferedFiles: make(map[string]int64),
		Receiver:        make(chan string),
		numOfMulti:      numOfMulti,
	}
}

func (m *Multiple) HandleReceive(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transferedFiles[name]++

	if m.transferedFiles[name] == m.numOfMulti {
		delete(m.transferedFiles, name)
		utils.MergeMultiFiles(name)
	}
}

func (m *Multiple) Run() {
	for name := range m.Receiver {
		m.HandleReceive(name)
	}
}
