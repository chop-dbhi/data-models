package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
)

const repoName = "https://github.com/chop-dbhi/data-models"

func pathExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func cloneRepo() {
	cmd := exec.Command("git", "clone", "--branch", repoBranch, repoName, repoDir)

	if logrus.GetLevel() == logrus.DebugLevel {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("problem cloning repo: %s", err)
	}
}

func pullRepo() {
	remote := fmt.Sprintf("origin/%s", repoBranch)
	cmd := exec.Command("git", "pull", repoDir, remote)

	if logrus.GetLevel() == logrus.DebugLevel {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("problem updating repo: %s", err)
	}
}

// updateRepo clones or updates the repo and returns true
// if an update occurred.
func updateRepo() {
	if !pathExists(repoDir) {
		cloneRepo()
	} else {
		pullRepo()
	}

	rebuildCache()
}

// pollRepo periodically checks the repo for updates.
func pollRepo() {
	// Check for updates every hour.
	t := time.NewTicker(interval)

	for {
		select {
		case <-t.C:
			updateRepo()
		}
	}
}
