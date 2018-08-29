package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/ssm"
	"k8s.io/api/core/v1"
)

type awsSessHandler struct {
	Sess *session.Session
}

func newAwsSessionHandler() *awsSessHandler {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("us-east-1")},
		Profile: "default",
	}))

	return &awsSessHandler{Sess: sess}
}

func (s *awsSessHandler) kmsDecrypt(file string) string {
	fmt.Println("Decrypting kubeconfig file:", file)

	// Create KMS service client
	svc := kms.New(s.Sess)

	// Encrypted data
	blob, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Got error reading file: ", err)
		os.Exit(1)
	}

	// Decrypt the data
	result, err := svc.Decrypt(&kms.DecryptInput{CiphertextBlob: blob})

	if err != nil {
		fmt.Println("Got error decrypting data: ", err)
		os.Exit(1)
	}

	content := []byte(result.Plaintext)
	tmpfile, err := ioutil.TempFile("/tmp", "kubeconfig")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()
}

func (s *awsSessHandler) ssmByPathToK8sEnvVar(path string) []v1.EnvVar {
	if path == "" {
		fmt.Println("No paramStore path provided")
		params := []v1.EnvVar{}
		return params
	}
	fmt.Println("paramStore path provided:", path)

	// Create a SSM client with additional configuration
	svc := ssm.New(s.Sess)
	var decrypt = true
	out, err := svc.GetParametersByPath(&ssm.GetParametersByPathInput{Path: &path, WithDecryption: &decrypt})
	if err != nil {
		fmt.Println("Got error getting params: ", err)
		os.Exit(1)
	}
	var result []v1.EnvVar
	for _, v := range out.Parameters {
		name := *v.Name
		value := *v.Value
		// Split the ssm path name
		sn := strings.Split(name, "/")
		// Get just the last string
		n := sn[len(sn)-1]
		// Set the params
		params := v1.EnvVar{Name: n, Value: value}
		result = append(result, params)
	}

	return result
}
