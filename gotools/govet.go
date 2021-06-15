package gotools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
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

	formatted, ok := formatVetInJSON(out.Bytes())
	if !ok {
		return nil, errors.New("failed to format vet output in JSON")
	}

	return formatted, nil
}

func formatVetInJSON(in []byte) ([]byte, bool) {
	objs := make([]string, 0)

	var obj string
	for _, s := range strings.Split(string(in), "\n") {
		if obj == "" && s != "{" {
			continue
		}

		obj += ("\t" + s)
		if s == "}" {
			objs = append(objs, obj)
			obj = ""
		} else {
			obj += "\n"
		}
	}

	if obj != "" {
		return nil, false
	}

	var out []byte
	if len(objs) == 0 {
		out = []byte("[]")
	} else {
		out = []byte("[\n" + strings.Join(objs, ",\n") + "\n]")
	}

	ok := json.Valid(out)
	if !ok {
		return nil, false
	}

	return out, true
}
