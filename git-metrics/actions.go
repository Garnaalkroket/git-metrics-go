package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Measurer interface {
	Measure()
	InitMeasurement()
}

type Action struct {
	method   string
	repo     string
	dir      string
	log_file string
}

type PullAction struct {
	Action
	key *ssh.PublicKeys
}

type CloneAction struct {
	Action
	key *ssh.PublicKeys
}

func (action Action) getAuthOptions() git.PullOptions {
	return git.PullOptions{}
}

func (action PullAction) Measure() {
	_, err := os.Stat(action.dir)
	if err != nil {
		publicKeys := action.key
		_, err := git.PlainClone(action.dir, false, &git.CloneOptions{
			URL:      action.repo,
			Progress: os.Stdout,
			Auth:     publicKeys,
		})
		CheckIfError(err)
	}
	publicKeys := action.key
	r, err := git.PlainOpen(action.dir)
	CheckIfError(err)
	w, err := r.Worktree()
	CheckIfError(err)
	start := time.Now()
	err = w.Pull(&git.PullOptions{RemoteName: "origin", Auth: publicKeys})
	CheckIfError(err)
	stop := time.Now()
	elapsed := stop.Sub(start)
	Info("Pull start " + start.String())
	Info("Pull end " + stop.String())
	Info("Pulling took " + elapsed.String())
	action.WriteToLogFile(elapsed)

}

func (action CloneAction) Measure() {
	start := time.Now()
	publicKeys := action.key
	_, err := git.PlainClone(action.dir, false, &git.CloneOptions{
		URL:      action.repo,
		Progress: os.Stdout,
		Auth:     publicKeys,
	})
	CheckIfError(err)
	stop := time.Now()
	elapsed := stop.Sub(start)
	Info("Clone start " + start.String())
	Info("Clone end " + stop.String())
	action.WriteToLogFile(elapsed)
	action.RemoveDir()

}

func (action Action) WriteToLogFile(elapsed time.Duration) {
	f, err := os.OpenFile(action.log_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%s RUN '%.2f' REPO '%s' \n", action.method, elapsed.Seconds(), action.repo))
	CheckIfError(err)
}

func (action Action) RemoveDir() {
	err := os.RemoveAll(action.dir)
	CheckIfError(err)
}

func (action Action) InitMeasurement() {
	_, err := os.Stat(action.dir)
	if err == nil {
		action.RemoveDir()
	}
}
