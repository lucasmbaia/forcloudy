package template

import (
  _template "text/template"
  "bytes"
  "fmt"
  "os"
)

func ConfGenerate(path, key, model string, conf interface{}) error {
  var (
    file *os.File
    err	  error
    t	  *_template.Template
    buf	  bytes.Buffer
  )

  if path != "" {
    if path[:1] != "/" {
      path += "/"
    }
  }

  if file, err = os.Create(fmt.Sprintf("%s%s", path, key)); err != nil {
    return err
  }

  defer func() {
    err = file.Close()
  }()

  t = _template.Must(_template.New("HA").Parse(model))

  if err = t.Execute(&buf, conf); err != nil {
    return err
  }

  if _, err = file.Write(buf.Bytes()); err != nil {
    return err
  }

  return err
}
