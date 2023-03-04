package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type TrackJSON struct {
	Id         int ``
	Start      string
	End        string
	Tags       []string
	Annotation string
}

type Worklog struct {
	Comment          string
	Started          string
	TimeSpentSeconds int
	Issue            string
}

func main() {
	var wg sync.WaitGroup
	var trackJSONRecords []TrackJSON
	var worklogs []Worklog

	err := json.NewDecoder(os.Stdin).Decode(&trackJSONRecords)
	if err != nil {
		log.Fatal(err)
	}

	for _, trackJSON := range trackJSONRecords {
		tagsLength := len(trackJSON.Tags)
		if tagsLength != 1 {
			log.Fatal("Invalid amount of tags")
		}
		start, err := time.Parse("20060102T150405Z", trackJSON.Start)
		if err != nil {
			log.Fatal(err)
		}
		end, err := time.Parse("20060102T150405Z", trackJSON.End)
		if err != nil {
			log.Fatal(err)
		}
		worklogs = append(worklogs, Worklog{
			Comment:          trackJSON.Annotation,
			Started:          start.Format("2006-01-02T15:04:05.000+0000"),
			TimeSpentSeconds: int(end.Sub(start).Seconds()),
			Issue:            trackJSON.Tags[0],
		})
	}

	for _, worklog := range worklogs {
		wg.Add(1)
		go issueWorklog(worklog, &wg)
	}

	wg.Wait()
}

func issueWorklog(worklog Worklog, wg *sync.WaitGroup) {
	issueMap := map[string]interface{}{
		"comment":          worklog.Comment,
		"started":          worklog.Started,
		"timeSpentSeconds": worklog.TimeSpentSeconds,
	}
	json_issue, err := json.Marshal(issueMap)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("https://jsonplaceholder.typicode.com/posts",
		"application/json",
		bytes.NewBuffer(json_issue))

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res)
	defer wg.Done()
}
