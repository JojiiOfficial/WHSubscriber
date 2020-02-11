package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/go-playground/webhooks.v5/github"
)

//GithubActionStruct the struct for github webhooks
type GithubActionStruct struct {
	action githubActionItem
}

type githubActionItem struct {
	on      github.Event
	actions []string
}

func createDefaultGithubFile(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		return nil
	}
	return ioutil.WriteFile(file, []byte("action:\n  on: push\n  filter:\n    branch:\n      master\n  actions:\n    - <call scripts>"), 0700)
}
