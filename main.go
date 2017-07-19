package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/briandowns/spinner"
	"gopkg.in/yaml.v2"
)

var (
	client = &http.Client{}
	token  string
	repo   string
)

// Label struct for GitHub
type Label struct {
	Name  string `json:"name" yaml:"name"`
	Color string `json:"color" yaml:"color"`
}

func loadLables(filename string) map[string]Label {
	_, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	var labels []Label
	err = yaml.Unmarshal(data, &labels)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	mappedLabels := make(map[string]Label)
	for _, l := range labels {
		mappedLabels[l.Name] = l
	}

	return mappedLabels
}

func getLabels() map[string]Label {
	url := fmt.Sprintf("https://api.github.com/repos/%s/labels", repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	var labels []Label
	err = json.Unmarshal(body, &labels)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	mappedLabels := make(map[string]Label)
	for _, l := range labels {
		mappedLabels[l.Name] = l
	}

	return mappedLabels
}

func updateLabel(label Label) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/labels/%s", repo, label.Name)

	body, err := json.Marshal(label)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer res.Body.Close()
}

func createLabel(label Label) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/labels", repo)

	body, err := json.Marshal(label)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer res.Body.Close()
}

func deleteLabel(label Label) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/labels/%s", repo, label.Name)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer res.Body.Close()
}

func updateLabels(userLabels map[string]Label, currentLabels map[string]Label) {
	for key, value := range userLabels {
		if _, ok := currentLabels[key]; ok {
			// Update Label
			delete(currentLabels, key)
			updateLabel(value)
		} else {
			// Create label
			createLabel(value)
		}
	}

	for _, value := range currentLabels {
		// Delete label
		deleteLabel(value)
	}
}

func main() {
	// Get Token
	fmt.Print("GitHub API Token: ")
	fmt.Scanln(&token)

	// Get Repo
	fmt.Print("GitHub Repo (owner/repo): ")
	fmt.Scanln(&repo)

	// File is always labels.yml for now
	file := "labels.yml"

	// Display spinner to indicate it's doing something
	s := spinner.New(spinner.CharSets[11], 50*time.Millisecond)
	s.Start()

	// Load user labels from file
	s.Suffix = " Loading labels from file"
	userLabels := loadLables(file)

	// Get repo labels from github
	s.Suffix = " Loading labels from GitHub repo"
	currentLabels := getLabels()

	// Update/Create/Delete labels based on the provided labels
	s.Suffix = " Updating labels on GitHub repo"
	updateLabels(userLabels, currentLabels)

	s.Stop()
	fmt.Printf("Labels Updated: https://github.com/%s/labels\n", repo)
}
