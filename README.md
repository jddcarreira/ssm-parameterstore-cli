SSM Parameter Store CLI
==========================
### Install
    Copy the binary for your architeture from bin/ to your $PATH

### Usage
```
Usage of ./ssm-manager:
  -aws-region string
    	AWS region to fetch the parameters. (default "eu-west-1")
  -key string
    	Key of the parameter stored.
  -op string
    	Operation to run. Options: put/get/get-all/del
  -silent
    	Silents any interaction and verification from the user
  -type string
    	Type of the parameter to store. Options: SecureString/String/StringList (default "SecureString")
  -value string
    	Value of the parameter to be stored.
```
#### Put parameter
    ./ssm-manager -op put -key /my/parameterstore/key -value blaplus2
#### Delete parameter
    ./ssm-manager -op del -key /my/parameterstore/key 
#### Get parameter
    ./ssm-manager -op get -key /my/parameterstore/key
#### Get all parameters
    ./ssm-manager -op get-all -key /my/parameterstore/key

