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
		MaxNumberOfMessages int64 `env:"SQS_MAX_MESSAGES"`
	}
}

func init() {

	if _, err := env.UnmarshalFromEnviron(&Env); err != nil {
		panic(err)
	}
}
