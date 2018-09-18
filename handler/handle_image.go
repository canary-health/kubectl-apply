package handler

import "os"

func setImage(imageFlag string) string {
	var image string
	if imageFlag == "" {
		image = os.Getenv("ECR_URL") + ":" + os.Getenv("TAG")
	}
	if image == ":" {
		panic("Image can not be empty try setting one using --image")
	}
	return image
}
