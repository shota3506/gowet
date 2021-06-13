package gotools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func GoListDir(path, dir string) (string, error) {
	cmd := exec.Command("go", "list", "-json", "-m", path)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run '%v': %w", cmd, err)
	}

	var resp struct {
		Dir string
	}
	if err := json.NewDecoder(&out).Decode(&resp); err != nil {
		return "", fmt.Errorf("failed to decode json: %w", err)
	}

	return resp.Dir, nil
}
