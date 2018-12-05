package template

import (
  "testing"
)

const model = `{{.Name}}`

type TestModel struct {
  Name	string	`json:",omitempty"`
}

func Test_ConfGenerate(t *testing.T) {
  var m = TestModel{Name: "MODEL"}

  if err := ConfGenerate("", "teste.cgf", model, m); err != nil {
    t.Fatalf("Erro to generate conf: %s", err.Error())
  }
}
