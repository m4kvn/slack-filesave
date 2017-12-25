package main

import (
	"github.com/nlopes/slack"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"io/ioutil"
	"os"
	"sync"
	"fmt"
	"flag"
	"time"
)

const folderName = "slack-downloads"

func main() {
	godotenv.Load()
	slackToken := *flag.String("token", os.Getenv("SLACK_API_TOKEN"), "Set Slack API Token")
	fileType := *flag.String("type", "image", "Set file type")
	doDelete := *flag.Bool("delete", false, "Deleting files from Slack")

	api := slack.New(slackToken)
	files, paging, err := getFiles(api, fileType, 1)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(folderName); err != nil {
		os.Mkdir(folderName, 0777)
	}

	deleteQue := make(chan string)
	waitGroup := sync.WaitGroup{}
	for paging.Page <= paging.Pages {
		for _, slackFile := range files {
			if !slackFile.IsPublic {
				continue
			}
			if _, err := os.Stat(getFileName(slackFile)); err == nil {
				log.Println("Already exists", slackFile.URLPrivateDownload)
				continue
			}
			waitGroup.Add(1)
			go write(&waitGroup, slackFile, slackToken)
		}
		log.Printf("files size: %d, paging: %#v\n", len(files), paging)
		files, paging, err = getFiles(api, fileType, paging.Page+1)
		if err != nil {
			log.Fatal(err)
		}
	}

	stopDelete := make(chan struct{})
	deleteWaitGroup := sync.WaitGroup{}
	if doDelete {
	deleteLoop:
		for {
			select {
			case fileId := <-deleteQue:
				deleteWaitGroup.Add(1)
				err := api.DeleteFile(fileId)
				if err != nil {
					log.Println(err)
				}
				deleteWaitGroup.Done()
				time.Sleep(1100)
			case <-stopDelete:
				break deleteLoop
			}
		}
	}
	waitGroup.Wait()
	close(stopDelete)
	deleteWaitGroup.Wait()
}

func write(waitGroup *sync.WaitGroup, slackFile slack.File, slackToken string) {
	defer waitGroup.Done()
	req, err := http.NewRequest("GET", slackFile.URLPrivateDownload, nil)
	req.Header.Set("Authorization", "Bearer "+slackToken)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	file, err := os.OpenFile(getFileName(slackFile), os.O_CREATE|os.O_WRONLY, 0666, )
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	_, err = file.Write(body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Download finished: %s\n", slackFile.URLPrivateDownload)
}

func getFiles(api *slack.Client, fileType string, page int) ([]slack.File, *slack.Paging, error) {
	return api.GetFiles(slack.GetFilesParameters{
		Types: fileType,
		Count: 1000,
		Page:  page,
	})
}

func getFileName(slackFile slack.File) string {
	return fmt.Sprintf("%s/%s-%s", folderName, slackFile.ID, slackFile.Name)
}
