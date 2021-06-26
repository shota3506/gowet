package gotools

import (
	"fmt"
	"os/exec"
)

func GoModInit(dir string) error {
	cmd := exec.Command("go", "mod", "init", "example.com")
	cmd.Dir = dir

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run 'go mod init': %w", err)
	}

	return nil
}
