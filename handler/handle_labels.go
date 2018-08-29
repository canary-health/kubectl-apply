package handler

import (
	"fmt"
	"strings"
)

func labelHandler(val string) map[string]string {
	ss := strings.Split(val, ",")
	out := make(map[string]string, len(ss))
	for _, pair := range ss {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 {
			fmt.Printf("%s must be formatted as key:value", pair)
		}
		out[kv[0]] = kv[1]
	}
	return out
}
