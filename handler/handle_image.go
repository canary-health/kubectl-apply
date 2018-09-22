package handler

import (
	"fmt"
	"os"
)

func setImage(imageFlag string) string {
	var image string
	if imageFlag == "" {
		image = os.Getenv("ECR_URI") + ":" + os.Getenv("TAG")
	}
	if image == ":" {
		panic("Image can not be empty try setting one using --image")
	}
	fmt.Println("Using image:", image)
	return image
}
