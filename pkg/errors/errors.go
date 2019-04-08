package errors

import (
	"github.com/autom8ter/objectify"
	"github.com/sfreiberg/gotwilio"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tool = objectify.New()

func Status(err error, code codes.Code, format, msg string) error {
	err = tool.WrapErrf(err, format, msg)
	return status.Error(code, err.Error())
}

func StatusStack(err error, code codes.Code, format, msg string) error {
	return StatusStack(err, code, format, msg)
}

func TwilioError(err error, code codes.Code, ex *gotwilio.Exception, msg string) error {
	return StatusStack(err, code, "Exception: %s\n", string(tool.MarshalJSON(ex)))
}
