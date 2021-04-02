package priceupdate

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sync"
)

// RunUpdate is GCP entry point for http trigger
func RunUpdate(w http.ResponseWriter, r *http.Request) {
	err := Run()
	if err != nil {
		panic(fmt.Sprintf("error: %s\n", err))
	}
	fmt.Fprintf(w, "Done")
}

func Run() error {
	env, err := getEnv()
	if err != nil {
		return err
	}

	// Populate from external resources
	securities, rows := MustPopulateSecuritiesAndRows(env["INPUT"], env["OUTPUT"])

	// existing output
	_output, err := rows.AsBytes()
	if err != nil {
		return err
	}

	// For each security create a task to retrieve a price from each source
	var mutex = sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(*securities))
	for i := 0; i < len(*securities); i++ {
		go GetISINSourcesPrice(&(*securities)[i], rows, &wg, &mutex)
	}
	wg.Wait()

	// new output
	output, err := rows.AsBytes()
	if err != nil {
		return err
	}

	// No change in output
	if bytes.Equal(_output, output) {
		return fmt.Errorf("output unchanged")
	}

	// Save output back to external store
	if err := rows.Save(env["API"], env["TOKEN"], output); err != nil {
		return fmt.Errorf("failed to save: %s", err)
	}

	return nil
}

func getEnv() (map[string]string, error) {
	env := make(map[string]string)
	for _, e := range []string{"INPUT", "OUTPUT", "API", "TOKEN"} {
		var v string
		v, ok := os.LookupEnv(e)
		if ok == false {
			return nil, fmt.Errorf("env not set: %s", e)
		}
		env[e] = v
	}
	return env, nil
}
