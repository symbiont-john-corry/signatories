package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/require"
)

var CREATOR string
var SIGNER string
var jobID string
var jobCompleteEvent Event

type Event struct {
	Index   int
	Channel string
	Type    string
	Data    ResultData
	JobID   string `json:"job_id"`
}
type EventResponse struct {
	Data struct {
		FirstIndex int     `json:"first_index"`
		LastIndex  int     `json:"last_index"`
		MaxIndex   int     `json:"max_index"`
		Events     []Event `json:"events"`
	}
}
type ResultData struct {
	Result map[string]interface{}
}

func (d ResultData) ToContract(c *Contract) error {
	m, err := json.Marshal(d.Result)
	if err != nil {
		return err
	}

	err = json.Unmarshal(m, c)
	if err != nil {
		return err
	}
	return nil
}

type Contract struct {
	Text    string
	Creator string
	Channel string
}

type JobResponse struct {
	Data struct {
		JobID string `json:"job_id"`
	} `json:"data"`
}

func Test_Contract(t *testing.T) {

	t.Run("reset network", func(t *testing.T) {
		c, b, e := exec.Command(`sym`, `network`, `reset`), new(strings.Builder), new(strings.Builder)
		c.Stdout = b
		c.Stderr = e
		err := c.Run()
		if err != nil {
			t.Log(err)
			t.Log(e.String())
		}

		t.Log(b.String())
	})

	t.Run("upgrade network to v10", func(t *testing.T) {
		c, b, e := exec.Command(`sym`, `network`, `upgrade`, `-v`, `10`), new(strings.Builder), new(strings.Builder)
		c.Stdout = b
		c.Stderr = e
		err := c.Run()
		if err != nil {
			t.Log(err)
			t.Log(e.String())
		}

		t.Log(b.String())
	})

	t.Run("publish contract", func(t *testing.T) {
		path, err := os.Getwd()
		require.NoError(t, err)
		t.Logf("publishing contracts: %s", path)

		c, b, e := exec.Command(`sym`, `network`, `publish`, `-d`, path), new(strings.Builder), new(strings.Builder)
		c.Stdout = b
		c.Stderr = e
		err = c.Run()
		if err != nil {
			t.Log(err)
			t.Log(e.String())
		}

		t.Log(b.String())
	})

	t.Run("generating key aliases", func(t *testing.T) {
		// create the key aliases for CREATOR AND SIGNER
		c, b := exec.Command(`sym`, `network`, `create-key-alias`), new(strings.Builder)
		c.Stdout = b
		c.Run()
		CREATOR = strings.TrimSpace(b.String())
		t.Logf("CREATOR: %s", CREATOR)

		c, b = exec.Command(`sym`, `network`, `create-key-alias`), new(strings.Builder)
		c.Stdout = b
		c.Run()
		SIGNER = strings.TrimSpace(b.String())
		t.Logf("SIGNER: %s", SIGNER)
	})

	t.Run("create contract", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodPost,
			"http://localhost:8888/api/v1/contracts/signatories/10-1.0.0/contract_create",
			bytes.NewBuffer([]byte(`{"text":"foo"}`)),
		)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Symbiont-Key-Alias", CREATOR)
		require.NoError(t, err)
		c := http.Client{}

		res, err := c.Do(req)
		require.NoError(t, err)
		// require.Equal(t, http.StatusCreated, res.StatusCode)

		t.Log("Response")
		b, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var data JobResponse

		err = json.Unmarshal(b, &data)
		require.NoError(t, err)
		jobID = data.Data.JobID
		t.Logf("Job ID: %s", jobID)
	})

	var contractChannel string

	t.Run("wait for contract data", func(t *testing.T) {
		o := func() error {
			err := getContractEventIndex(jobID)
			if err != nil {
				t.Logf("error getting event: %v", err)
			}
			return err
		}
		err := backoff.Retry(o, backoff.NewExponentialBackOff())
		require.NoError(t, err)

		t.Logf("Job Complete Event: %+v", jobCompleteEvent)

		var c Contract
		err = jobCompleteEvent.Data.ToContract(&c)
		require.NoError(t, err)
		contractChannel = string(c.Channel)
	})

	var contract Contract

	t.Run("get contract", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodPost,
			"http://localhost:8888/api/v1/contracts/signatories/10-1.0.0/contract",
			bytes.NewBuffer([]byte(fmt.Sprintf(`{"channel":"%s"}`, contractChannel))),
		)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Symbiont-Key-Alias", CREATOR)
		require.NoError(t, err)
		c := http.Client{}

		res, err := c.Do(req)
		require.NoError(t, err)

		b, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		var e Event
		err = json.Unmarshal(b, &e)
		require.NoError(t, err)

		err = e.Data.ToContract(&contract)
		require.NoError(t, err)
	})

	t.Run("sign contract - creator", func(t *testing.T) {

	})

	t.Run("sign contract - other party", func(t *testing.T) {

	})

	t.Run("add party to contract", func(t *testing.T) {

	})

	t.Run("sign contract - other party", func(t *testing.T) {

	})
}

func getContractEventIndex(jid string) error {
	eventURL := "http://localhost:8888/api/v1/events/foo"

	req, err := http.NewRequest(
		http.MethodGet,
		eventURL,
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Symbiont-Key-Alias", CREATOR)
	c := http.Client{}

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 status: %d; %s", res.StatusCode, string(b))
	}

	var er EventResponse
	err = json.Unmarshal(b, &er)
	if err != nil {
		return err
	}

	for _, e := range er.Data.Events {
		if e.JobID == jid {
			if e.Type == `assembly/job_complete` {
				jobCompleteEvent = e
				return nil
			}
		}
	}

	return fmt.Errorf("no completed job found")
}
