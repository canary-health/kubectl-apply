package handler

import (
	"fmt"
	"io/ioutil"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

type decoderHandler struct {
	Obj runtime.Object
}

func newDecoderHandler(file string) *decoderHandler {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		panic(err)
	}
	dat, err := ioutil.ReadFile(file)
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(dat), nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}
	return &decoderHandler{obj}
}
