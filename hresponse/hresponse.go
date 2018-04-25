package hresponse

import "github.com/sirupsen/logrus"

type HresponseConfig struct {
	Logger logrus.FieldLogger
}

var hresponse HresponseConfig

func New (config HresponseConfig){
	hresponse = config
}