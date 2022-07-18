package jobs

import (
	"os"
	"strconv"
)

type Multiple struct {
	transferedFiles map[string]int64
	receiver        chan string
	numOfMulti      int64
}

func MultipleJob() *Multiple {
	numOfMulti, _ := strconv.ParseInt(os.Getenv("NUM_OF_MULTI"), 10, 64)
	return &Multiple{
		transferedFiles: make(map[string]int64),
		receiver:        make(chan string),
		numOfMulti:      numOfMulti,
	}
}

func (m *Multiple) Run() {
	for message := range m.receiver {
		m.transferedFiles[message]++
		if m.transferedFiles[message] == m.numOfMulti {
		}
	}
}
