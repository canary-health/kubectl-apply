package handler

import (
	"fmt"
	"os"
	"strings"
)

func labelHandler(labelsFlag string) map[string]string {
	if labelsFlag == "" && os.Getenv("CODEBUILD_RESOLVED_SOURCE_VERSION") != "" {
		out := make(map[string]string, 1)
		value := fmt.Sprintf("%v", os.Getenv("CODEBUILD_RESOLVED_SOURCE_VERSION")[0:8])
		out["build"] = value
		return out
	}
	if labelsFlag == "" && os.Getenv("CODEBUILD_RESOLVED_SOURCE_VERSION") == "" {
		out := make(map[string]string, 1)
		out["build"] = "unavailable"
		return out
	}
	ss := strings.Split(labelsFlag, ",")
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
