package aws

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

// Secret struct
type Secret struct {
	Key        string `json:"key"`
	ID         string `json:"id"`
	ExternalID string `json:"externalID"`
}

// GetSecret return secret for logining prisma api
func GetSecret(secretName string, svc secretsmanageriface.SecretsManagerAPI) (*Secret, error) {
	if secretName == "" {
		return nil, errors.New("Parameter secretName is missing")
	}

	result, err := svc.GetSecretValue(
		&secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(secretName),
			VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
		})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return nil, err
	}
	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			return nil, err
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
	}

	secret := &Secret{}
	if secretString != "" {
		json.Unmarshal([]byte(secretString), &secret)
	} else if decodedBinarySecret != "" {
		json.Unmarshal([]byte(decodedBinarySecret), &secret)
	}

	return secret, nil
}
