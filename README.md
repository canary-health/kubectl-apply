Encrypt kube configs following this: https://docs.aws.amazon.com/cli/latest/reference/kms/encrypt.html

```
aws kms encrypt --key-id <AWS_US-EAST-1_KEY-ID> --plaintext fileb://<PATH_TO_KUBECONFIG_FILE> --output text --query CiphertextBlob | base64 --decode > ./config/kube/<OUTPUT_ENCRYPTED_FILE_NAME>
```

Example `kubectl-apply` commands:

Kind: Deployment

```
kubectl-apply deployment --kubeconfig="../config/kube/config-aws-canary-dev" -f ./test/testdata/deployment.stage.yml --labels="test:true,another:true"
```

Kind: Service

```
kubectl-apply service --kubeconfig="../config/kube/config-aws-canary-dev" -f ./test/testdata/service.stage.yml --labels="test:true,another:true"
```

Kind: CronJob

```
kubectl-apply cronjob --kubeconfig="../config/kube/config-aws-canary-dev" -f ./test/testdata/cronjob.yml --labels="test:true,another:true"
```
