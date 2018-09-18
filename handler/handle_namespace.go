package handler

import "os"

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
	return namespace
}
