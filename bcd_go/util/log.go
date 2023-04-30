package util

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var Log = initLog()

func initLog() *zap.SugaredLogger {
	temp, err := zap.NewDevelopment()
	if err != nil {
		panic("log init fail:" + err.Error())
	}
	return temp.Sugar()
}

func LogErrorStack(err error) {
	if err != nil {
		Log.Errorf("%+v", err)
	}
}

func LogErrorsStack(errs []error, msg string, args ...interface{}) error {
	for _, e := range errs {
		Log.Errorf("%+v", e)
	}
	return errors.Errorf(msg, args...)
}
