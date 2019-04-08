package goconnect_test

import (
	"github.com/autom8ter/goconnect"
	"github.com/autom8ter/goconnect/plugins"
	"testing"
)

func Test(t *testing.T) {
	if err := goconnect.NewFromFileEnv("credentials.json").Serve(":3000", plugins.EchoService()); err != nil {
		t.Fatal(err.Error())
	}
}
