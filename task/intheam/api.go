package intheam

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type API struct {
	client *http.Client
}

func New() *API {
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	return &API{
		client: &http.Client{Transport: tr},
	}
}

func (api *API) GetTriggers(apiKey string) []Task {
	req, err := http.NewRequest("GET", "https://inthe.am/api/v2/tasks/", nil)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	req.Header.Add("Authorization", "TOKEN "+apiKey)
	res, err := api.client.Do(req)
	if err != nil {
		log.Printf("err: %v", err)
	}
	// TODO: check status and return error

	tasksRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()

	var tasks []Task
	err = json.Unmarshal(tasksRaw, &tasks)
	if err != nil {
		log.Printf("err: %v\nrecv: %v", err, string(tasksRaw))
	}
	// TODO: handle error
	return tasks
}

func (api *API) NewTrigger(apiKey, description string) {
	task := Task{
		Description: description,
		Tags:        []string{"review"},
		Status:      "pending",
	}
	taskJSON, _ := json.Marshal(task)
	req, err := http.NewRequest("POST", "https://inthe.am/api/v2/tasks/", bytes.NewBuffer(taskJSON))
	if err != nil {
		log.Printf("err: %v", err)
	}
	req.Header.Add("Authorization", "TOKEN "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := api.client.Do(req)
	if err != nil {
		log.Printf("err: %v", err)
	}
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("err : %v", err)
	}
	res.Body.Close()

	log.Printf("resp: %v", string(respBody))
}
