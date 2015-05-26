package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	updatingRepo = false
	updateLock   = sync.Mutex{}
)

func pathExists(p string) bool {
	_, err := os.Stat(p)

	return err == nil
}

func cloneRepo() {
	cmd := exec.Command("git", "clone", "--branch", repoBranch, repoName, repoDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("problem cloning repo: %s", err)
	}
}

func pullRepo() {
	remote := fmt.Sprintf("origin/%s", repoBranch)
	cmd := exec.Command("git", "-C", repoDir, "pull", ".", remote)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("problem pulling repo: %s", err)
	}
}

// updateRepo clones or updates the repo and returns true
// if an update occurred.
func updateRepo() {
	// Update already in progress
	if updatingRepo {
		return
	}

	updateLock.Lock()
	defer updateLock.Unlock()

	updatingRepo = true

	gitDir := filepath.Join(repoDir, ".git")

	if !pathExists(gitDir) {
		cloneRepo()
	} else {
		pullRepo()
	}

	rebuildCache(repoDir)

	updatingRepo = false
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
