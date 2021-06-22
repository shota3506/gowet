package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func marshalVet(in []byte) ([]byte, error) {
	objs := make([]string, 0)

	flag := false
	obj := ""
	for _, s := range strings.Split(string(in), "\n") {
		if !flag && s != "{" {
			continue
		}

		if s == "{" {
			flag = true
		} else if s == "}" {
			objs = append(objs, obj)
			obj = ""
			flag = false
		} else {
			if obj != "" {
				obj += "\n"
			}
			obj += s
		}
	}

	if flag || obj != "" {
		return nil, errors.New("failed to format vet output in JSON")
	}

	var out []byte
	if len(objs) == 0 {
		out = []byte("{}")
	} else {
		out = []byte("{\n" + strings.Join(objs, ",\n") + "\n}")
	}

	if !json.Valid(out) {
		return nil, errors.New("failed to format vet output in JSON: invalid JSON")
	}
	var bf bytes.Buffer
	if err := json.Indent(&bf, out, "", "\t"); err != nil {
		return nil, fmt.Errorf("failed to format vet output in JSON: %w", err)
	}

	return bf.Bytes(), nil
}
