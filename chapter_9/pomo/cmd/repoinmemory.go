package cmd

import (
	"rohitsingh/pomo/pomodoro"
	"rohitsingh/pomo/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
