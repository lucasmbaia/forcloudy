package template

import (
  _template "text/template"
  "bytes"
  "fmt"
  "os"
)

const (
  MINION_SERVER = `
{{$name := .Name}}
{{$dns := .Dns}}
{{range .Hosts -}}
{{printf "frontend f_%s-%s" $name .PortSRC}}
{{if ne .PortSRC "443" -}}
{{printf "\tbind *:443 ssl crt /etc/haproxy/ssl/totino.com.br.pem crt /etc/haproxy/ssl/totino2.com.br.pem"}}
{{printf "\tmode http"}}
{{printf "\tacl %s_acl hdr(host) -i %s" $name $dns}}
{{printf "\tuse_backend b_%s-%s if %s_acl" $name .PortSRC $name}}
{{printf "\n"}}
{{printf "backend b_%s-%s" $name .PortSRC}}
{{printf "\tmode http"}}
{{printf "\tserver minion-1 minion-1.com.br:443 check ssl verify none"}}
{{printf "\tlog /dev/log local10 debug"}}
{{else -}}
{{printf "\tbind *:%s" .PortSRC}}
{{printf "\tmode tcp"}}
{{printf "\tuse_backend b_%s-%s" $name .PortSRC}}
{{printf "\n"}}
{{printf "backend b_%s-%s" $name .PortSRC}}
{{printf "\tmode tcp"}}
{{printf "\tserver minion-1 minion-1.com.br:%s check" .PortSRC}}
{{end}}
{{end}}`
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
