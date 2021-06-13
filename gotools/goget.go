package gotools

import (
	"fmt"
	"os/exec"
)

func GoGet(path, dir string) error {
	cmd := exec.Command("go", "get", path)
	cmd.Dir = dir

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run '%v': %w", cmd, err)
	}

	return nil
}
