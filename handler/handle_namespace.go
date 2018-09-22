package handler

import (
	"fmt"

	"github.com/canary-health/gotils/pkg/env"
)

func setNamespace(nsFlag string) string {
	var namespace string
	if nsFlag == "" {
		// if passed flag not set try the var $TAG from build
		namespace = env.Get("TAG", "development")
	} else {
		namespace = nsFlag
	}
	fmt.Println("Setting namespace:", namespace)
	return namespace
}
