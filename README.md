[![Go Report Card](https://goreportcard.com/badge/github.com/schrius/prisma-cloud-hub)](https://goreportcard.com/report/github.com/schrius/prisma-cloud-hub)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![Go](https://github.com/schrius/prisma-cloud-hub/workflows/Go/badge.svg)

# Prisma Cloud Remediation
Prisma cloud remediation is an application build with AWS SAM.
A SQS is integrated with Prisma to receive alert, and distribute
to a corresponding lambda function to resolve the issue.

## VPC Terminator
AWS Lambda Function for cleanup VPC Resouce in Non Virginia Region. 
Required integration with Prisma Cloud, with Policy Query:
```
config where cloud.type = 'aws' AND cloud.region != 'AWS Virginia' AND api.name = 'aws-ec2-describe-vpcs' AND json.rule = state equals "available"
```
Alert Use AWS SQS to triger clean up.

## Suspicious Traffic
Send SMS to Cloud Services team when suspicious traffic was detected by Prisma

## False Alert Remover
Verify Prisma alert and dismiss false alert in Prisma

## prisma-aws-onboard
A Lambda function use to register new account onboarding on Prisma.

## Dispatcher
The entry point of the remediation. See [link](./docs/prisma_integration_architecture.md) for more details.

### File
```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── dispatcher                   <-- Source code for dispatcher lambda function
│   ├── main.go                 <-- Lambda function code
│   └── main_test.go            <-- Unit tests
└── template.yaml
```
### Install SAM
##### Mac or Linux
Install SAM CLI using Brew

brew tap aws/tap
brew install aws-sam-cli

##### Windows
[Download 64 bit](https://github.com/awslabs/aws-sam-cli/releases/latest/download/AWS_SAM_CLI_64_PY3.msi)

### AWS SAM Reference
[Reference](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-command-reference.html)

### Build 
#### Using Makefile
```bash
make
```
The Makefiles contains all necessary steps to install dependencies, build, static analysis, and test all code.

#### SAM Build
```bash
sam build
```
SAM all install all depenecies and build the artificate but it does not perform unit test.

#### Deployment
To deploy the application
```bash
sam deploy
```
SAM will package all artifact, upload to S3 bucket and deploy the application
## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)

## Go

### Installing dependencies
```shell
go get -u ./...
```
### Testing
```shell
go test -v ./...
```

### Building

Golang is a statically compiled language, meaning that in order to run it you have to build the executable target.

You can issue the following command in a shell to build it:

```shell
GOOS=linux GOARCH=amd64 go build -o dispatcher/dispatcher ./dispatcher
```

**NOTE**: If you're not building the function on a Linux machine, you will need to specify the `GOOS` and `GOARCH` environment variables, this allows Golang to build your function for another system architecture and ensure compatibility.

```bash
sam local start-api
```

## Packaging and deployment

To deploy your application for the first time, run the following in your shell:

```bash
sam deploy
```

The command will package and deploy your application to AWS, with a series of prompts:

* **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be something matching your project name.
* **AWS Region**: The AWS region you want to deploy your app to.
* **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review. If set to no, the AWS SAM CLI will automatically deploy application changes.
* **Allow SAM CLI IAM role creation**: Many AWS SAM templates, including this example, create AWS IAM roles required for the AWS Lambda function(s) included to access AWS services. By default, these are scoped down to minimum required permissions. To deploy an AWS CloudFormation stack which creates or modified IAM roles, the `CAPABILITY_IAM` value for `capabilities` must be provided. If permission isn't provided through this prompt, to deploy this example you must explicitly pass `--capabilities CAPABILITY_IAM` to the `sam deploy` command.
* **Save arguments to samconfig.toml**: If set to yes, your choices will be saved to a configuration file inside the project, so that in the future you can just re-run `sam deploy` without parameters to deploy changes to your application.

You can find your API Gateway Endpoint URL in the output values displayed after deployment.

# Appendix

### Golang installation

Please ensure Go 1.x (where 'x' is the latest version) is installed as per the instructions on the official golang website: https://golang.org/doc/install

A quickstart way would be to use Homebrew, chocolatey or your linux package manager.

#### Homebrew (Mac)

Issue the following command from the terminal:

```shell
brew install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
brew update
brew upgrade golang
```

#### Chocolatey (Windows)

Issue the following command from the powershell:

```shell
choco install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
choco upgrade golang
```

## Bringing to the next level

Here are a few ideas that you can use to get more acquainted as to how this overall process works:

* Create an additional API resource (e.g. /hello/{proxy+}) and return the name requested through this new path
* Update unit test to capture that
* Package & Deploy

Next, you can use the following resources to know more about beyond hello world samples and how others structure their Serverless applications:

* [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/)
