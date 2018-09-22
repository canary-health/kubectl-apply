package handler

import (
	"fmt"
	"os"
)

func setNamespace(nsFlag string) string {
	var namespace string
	if nsFlag == "" {
		// if passed flag not set try the var $TAG from build
		namespace = os.Getenv("TAG")
	}
	// if $TAG also not set default to development
	if namespace == "" {
		namespace = "development"
	}
	fmt.Println("Setting namespace:", namespace)
	return namespace
}
