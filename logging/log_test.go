package logging

import (
	"testing"
)

const (
	TEMPLATE = "S=%s, SS=%s"
)

type S struct {
	Teste string
}

func Test_Infof(t *testing.T) {
	var l = New(INFO)
	var s = S{Teste: "lucas"}
	var p = &S{Teste: "martins"}

	l.Infofc(TEMPLATE, s, p)
}
