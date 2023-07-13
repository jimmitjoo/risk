-include .env
export

.PHONY: deps clean build

SourceDir=src
StackName=risk-game-stack
Profile=${AWS_PROFILE}

local:
	go run ${SourceDir}/router.go

build:
	GOOS=linux GOARCH=amd64 go build -o bin/router ${SourceDir}/router.go

deploy: build
	sam validate
	sam package --template-file template.yaml --s3-bucket risk-binar-bucket --output-template-file packaged.yaml
	aws cloudformation deploy --template-file packaged.yaml --stack-name ${StackName} --profile ${Profile} --capabilities CAPABILITY_NAMED_IAM

delete:
	aws cloudformation delete-stack --stack-name ${StackName} --profile ${Profile}