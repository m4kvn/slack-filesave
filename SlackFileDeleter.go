package main

import (
	"github.com/nlopes/slack"
	"sync"
	"time"
	"log"
)

type SlackFileDeleter struct {
	DeleteQue       chan string
	StopChan        chan struct{}
	deleteWaitGroup sync.WaitGroup
}

func newInstance() SlackFileDeleter {
	return SlackFileDeleter{
		DeleteQue:       make(chan string, 1000),
		StopChan:        make(chan struct{}),
		deleteWaitGroup: sync.WaitGroup{},
	}
}

func (s *SlackFileDeleter) delete(fileId string) {
	s.deleteWaitGroup.Add(1)
	go func() {
		s.DeleteQue <- fileId
	}()
}

func (s *SlackFileDeleter) run(api *slack.Client) {
	for {
		select {
		case fileId := <-s.DeleteQue:
			if err := api.DeleteFile(fileId); err != nil {
				log.Println(err)
			}
			log.Printf("Delete %s from Slack.\n", fileId)
			s.deleteWaitGroup.Done()
			time.Sleep(1 * time.Second)
		case <-s.StopChan:
			close(s.DeleteQue)
			return
		}
	}
}

func (s *SlackFileDeleter) stop() {
	s.deleteWaitGroup.Wait()
	close(s.StopChan)
}
