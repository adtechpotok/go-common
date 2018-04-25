package hresponse

import "github.com/adtechpotok/silog/logrus"

type HresponseConfig struct {
	Logger logrus.FieldLogger
}

var hresponse HresponseConfig

func New(config HresponseConfig) {
	hresponse = config
}
