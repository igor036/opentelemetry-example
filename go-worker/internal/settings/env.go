package settings

import (
	"github.com/Netflix/go-env"
)

var Env enviroment

type enviroment struct {
	ViaCep struct {
		BaseURL string `env:"VIA_CEP_BASE_URL"`
	}
	NodeAPI struct {
		BaseURL string `env:"NODE_API_BASE_URL"`
	}
	SQS struct {
		MaxNumberOfMessages int64  `env:"SQS_MAX_MESSAGES"`
		AWSRegion           string `env:"AWS_REGION"`
		AWSAccessKeyId      string `env:"AWS_ACCESS_KEY_ID"`
		AWSSecretAccessKey  string `env:"AWS_SECRET_ACCESS_KEY"`
		AWSAddress          string `env:"AWS_ADDRESS"`
		AWSProfile          string `env:"AWS_PROFILE"`
		ImportZipCodeQueue  string `env:"IMPORT_ZIP_CODE_QUEUE"`
	}
}

func init() {

	if _, err := env.UnmarshalFromEnviron(&Env); err != nil {
		panic(err)
	}
}
