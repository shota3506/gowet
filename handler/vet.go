package handler

import (
	"encoding/json"
	"strings"
)

func marshalVet(in []byte) ([]byte, bool) {
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

	if !json.Valid(out) {
		return nil, false
	}

	return out, true
}
