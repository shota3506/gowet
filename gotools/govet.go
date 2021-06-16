package gotools

import (
	"bytes"
	"fmt"
	"os/exec"
)

func GoVet(dir string) ([]byte, error) {
	cmd := exec.Command("go", "vet", "-json", "./...")
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run '%v': %w", cmd, err)
	}

	return out.Bytes(), nil
}
