package hresponse

import (
	"net/http"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/adtechpotok/silog"
)

func ResultJson(result interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if _, ok := result.(ErrorResponse); ok {
		w.WriteHeader(result.(ErrorResponse).status)
	}
	js, err := jsoniter.Marshal(result)
	if err != nil {
		silog.WithFields(logrus.Fields{"error": err,}).Warning("Json encode failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
