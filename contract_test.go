package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var CREATOR string
var SIGNER string
var contractID string

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
		t.Logf(string(b))
	})

	t.Run("get contract", func(t *testing.T) {
		_, _ = http.NewRequest(
			http.MethodPost,
			"http://localhost:8888/api/v1/contracts/signatories/10-1.0.0/contract",
			bytes.NewBuffer([]byte("foo")),
		)
	})
}
