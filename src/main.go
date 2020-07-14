package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func getSSMService(region string) (*ssm.SSM, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, fmt.Errorf("Cannot authenticate with AWS because: %s", err)

	}
	return ssm.New(sess), nil
}

func getAllParameters(ssmKey string, ssmConnection *ssm.SSM) error {
	var nextToken string
	var ssmKeys []string
	var ssmResults *ssm.GetParametersByPathOutput
	var err error

	for {
		if nextToken == "" {
			ssmResults, err = ssmConnection.GetParametersByPath(
				&ssm.GetParametersByPathInput{
					Path:           aws.String(ssmKey),
					WithDecryption: aws.Bool(true),
					Recursive:      aws.Bool(false),
				})
		} else {
			ssmResults, err = ssmConnection.GetParametersByPath(
				&ssm.GetParametersByPathInput{
					Path:           aws.String(ssmKey),
					WithDecryption: aws.Bool(true),
					Recursive:      aws.Bool(false),
					NextToken:      aws.String(nextToken),
				})

		}

		if err != nil {
			return fmt.Errorf("Cannot get parameters list because: %s", err)
		}

		for i := range ssmResults.Parameters {
			ssmKeys = append(ssmKeys, aws.StringValue(ssmResults.Parameters[i].Name))
		}

		if ssmResults.NextToken == nil {
			break
		}
		nextToken = aws.StringValue(ssmResults.NextToken)
	}

	for i := range ssmKeys {
		fmt.Println(ssmKeys[i])
	}
	return nil
}

func getParameter(ssmKey string, ssmConnection *ssm.SSM) error {
	ssmVal, err := ssmConnection.GetParameter(
		&ssm.GetParameterInput{
			Name:           &ssmKey,
			WithDecryption: aws.Bool(true),
		})

	if err != nil {
		return fmt.Errorf("Cannot get parameter because: %s", err)
	}

	fmt.Printf("Name: %s\nType: %s\nValue: %s\nVersion: %c\nLastModifiedDate: %s\nARN: %s",
		aws.StringValue(ssmVal.Parameter.Name),
		aws.StringValue(ssmVal.Parameter.Type),
		aws.StringValue(ssmVal.Parameter.Value),
		aws.Int64Value(ssmVal.Parameter.Version),
		aws.TimeValue(ssmVal.Parameter.LastModifiedDate),
		aws.StringValue(ssmVal.Parameter.ARN),
	)
	return nil
}

func putParameter(ssmKey string, ssmValue string, ssmValueType string, ssmConnection *ssm.SSM) error {
	parameter, err := ssmConnection.PutParameter(
		&ssm.PutParameterInput{
			Name:      aws.String(ssmKey),
			Type:      aws.String(ssmValueType),
			Value:     aws.String(ssmValue),
			Overwrite: aws.Bool(true),
		})

	if err != nil {
		return err
	}

	fmt.Printf("Parameter Saved!\nName: %s\nVersion: %d", ssmKey, aws.Int64Value(parameter.Version))
	return nil
}

func deleteParameter(ssmKey string, silent bool, ssmConnection *ssm.SSM) error {
	_, err := ssmConnection.DeleteParameter(
		&ssm.DeleteParameterInput{
			Name: aws.String(ssmKey),
		})

	if err != nil {
		return fmt.Errorf("Cannot delete parameter because: %s", err)
	}

	fmt.Printf("Parameter Deleted!\nName: %s", ssmKey)
	return nil
}

func main() {
	var err error
	var silent bool
	var awsRegion, op, key, value, valueType string
	flag.StringVar(&awsRegion, "aws-region", "eu-west-1", "AWS region to fetch the parameters.")
	flag.StringVar(&op, "op", "", "Operation to run. Options: put/get/get-all/del")
	flag.StringVar(&key, "key", "", "Key of the parameter stored.")
	flag.StringVar(&value, "value", "", "Value of the parameter to be stored.")
	flag.StringVar(&valueType, "type", "SecureString", "Type of the parameter to store. Options: SecureString/String/StringList")
	flag.BoolVar(&silent, "silent", false, "Silents any interaction and verification from the user")
	flag.Parse()

	if key == "" {
		fmt.Printf("Key or type was not provided")
		os.Exit(2)
	}

	ssmCon, err := getSSMService(awsRegion)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	switch op {
	case "get":
		err = getParameter(key, ssmCon)
	case "put":
		if value == "" {
			fmt.Printf("Value or type was not provided")
			os.Exit(2)
		}
		err = putParameter(key, value, valueType, ssmCon)
	case "get-all":
		err = getAllParameters(key, ssmCon)
	case "del":
		err = deleteParameter(key, silent, ssmCon)
	default:
		flag.PrintDefaults()
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
