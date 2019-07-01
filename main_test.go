package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hermesespinola/proxy-app/api/handlers"
	"github.com/hermesespinola/proxy-app/api/middlewares"
	"github.com/hermesespinola/proxy-app/api/server"
	"github.com/hermesespinola/proxy-app/utils"
)

func init() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		utils.LoadEnv()
		app := server.Setup()
		handlers.HandlerRedirection(app)
		wg.Done()
		server.RunServer(app)
	}(wg)
	wg.Wait()
}

type Response struct {
	Status string `json:"status,omitempty"`
	Queue  string `json:"queue,omitempty"`
	New    string `json:"new,omitempty"`
	Popped string `json:"popped,omitempty"`
}

type TestCase struct {
	Domain string
	Output string
}

func TestRead(t *testing.T) {
	client := http.Client{}
	cases := []TestCase{
		{
			Domain: "alpha",
			Output: "[{\"domain\":\"alpha\",\"weight\":5,\"priority\":5},{\"domain\":\"omega\",\"weight\":1,\"priority\":5},{\"domain\":\"beta\",\"weight\":5,\"priority\":1}]",
		},
	}

	for _, caze := range cases {
		req, err := http.NewRequest("GET", "http://localhost:8081/read", nil)
		req.Header.Add("domain", caze.Domain)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		bytes, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		res := Response{}
		json.Unmarshal(bytes, &res)
		assert.Equal(t, "ok", res.Status)
		assert.Equal(t, caze.Output, res.Queue)
	}
}

func TestAlgorithm(t *testing.T) {
	client := http.Client{}
	cases := []TestCase{
		{Domain: "alpha"},
		{Domain: "alpha"},
		{Domain: "alpha"},
	}

	for _, caze := range cases {
		req, err := http.NewRequest("GET", "http://localhost:8081/push", nil)
		req.Header.Add("domain", caze.Domain)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		bytes, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		res := Response{}
		json.Unmarshal(bytes, &res)
		assert.Equal(t, "ok", res.Status)
	}

	minPrior := math.MaxFloat64
	for range cases {
		req, err := http.NewRequest("GET", "http://localhost:8081/pop", nil)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		bytes, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		res := Response{}
		json.Unmarshal(bytes, &res)
		popped := middlewares.RepoNode{}
		json.Unmarshal([]byte(res.Popped), &popped)
		val := popped.Value()
		minPrior = math.Min(minPrior, val)
		assert.Equal(t, val, minPrior)
	}
}
