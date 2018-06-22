package template

const (
  MINION  = `
{{"global"}}
{{printf "\tstats socket /run/haproxy/admin.sock mode 777 level admin expose-fd listeners"}}
{{printf "\tstats timeout 30s"}}
{{printf "\tdaemon"}}
{{printf "\tmaxconn 2000"}}
{{printf "\n"}}
{{"defaults"}}
{{printf "\tlog \tglobal"}}
{{printf "\tmode \thttp"}}
{{printf "\tretries \t3"}}
{{printf "\toption http-keep-alive"}}
{{printf "\toption dontlognull"}}
{{printf "\ttimeout connect 5000"}}
{{printf "\ttimeout client 50000"}}
{{printf "\ttimeout server 50000"}}
{{$name := .Name}}
{{$dns := .Dns}}
{{range .Hosts -}}
{{printf "frontend f_%s-%s" $name .PortSRC}}
{{if eq .Protocol "https" -}}
{{printf "\tbind *:%s ssl crt /etc/haproxy/ssl/totino.com.br.pem crt /etc/haproxy/ssl/totino2.com.br.pem" .PortSRC}}
{{printf "\tmode http"}}
{{printf "\tacl %s_acl_https hdr(host) -i %s" $name $dns}}
{{printf "\tuse_backend b_%s-%s if %s_acl_https" $name .PortSRC $name}}
{{printf "\n"}}
{{printf "backend b_%s-%s" $name .PortSRC}}
{{printf "\tmode http"}}
{{range $idx, $addr := .Address -}}
{{printf "\tserver application-%d %s check" $idx $addr}}
{{end -}}
{{printf "\tlog /dev/log local0 debug"}}
{{printf "\n"}}

{{else -}}{{if eq .Protocol "http" -}}
{{printf "\tbind *:%s" .PortSRC}}
{{printf "\tmode http"}}
{{printf "\tlog /dev/log local0 debug"}}
{{printf "\tacl whitelist src %s" .Whitelist}}
{{printf "\tuse_backend b_%s-%s if whitelist" $name .PortSRC}}
{{printf "\n"}}

{{printf "backend b_%s-%s" $name .PortSRC}}
{{printf "\tmode http"}}
{{printf "\thttp-request set-header Host %s" $dns}}
{{range $idx, $addr := .Address -}}
{{printf "\tserver application-%d %s check" $idx $addr}}
{{end -}}
{{printf "\n"}}

{{else -}}
{{printf "\tbind *:%s" .PortSRC}}
{{printf "\tmode tcp"}}
{{printf "\tuse_backend b_%s-%s" $name .PortSRC}}
{{printf "\n"}}

{{printf "backend b_%s-%s" $name .PortSRC}}
{{printf "\tmode tcp"}}
{{range $idx, $addr := .Address -}}
{{printf "\tserver application-%d %s check" $idx $addr}}
{{end -}}
{{printf "\tsource 0.0.0.0 usesrc client"}}

{{end -}}
{{end -}}
{{end -}}`

  MINION_SERVER = `
{{"global"}}
{{printf "\tlog /dev/log \tlocal0"}}
{{printf "\tlog /dev/log \tlocal1 debug"}}
{{printf "\tchroot /var/lib/haproxy"}}
{{printf "\tstats socket /run/haproxy/admin.sock mode 660 level admin"}}
{{printf "\tstats timeout 30s"}}
{{printf "\tuser haproxy"}}
{{printf "\tgroup haproxy"}}
{{printf "\tdaemon"}}
{{printf "\tca-base /etc/ssl/certs"}}
{{printf "\tcrt-base /etc/ssl/private"}}
{{printf "\tssl-default-bind-ciphers ECDH+AESGCM:DH+AESGCM:ECDH+AES256:DH+AES256:ECDH+AES128:DH+AES:ECDH+3DES:DH+3DES:RSA+AESGCM:RSA+AES:RSA+3DES:!aNULL:!MD5:!DSS"}}
{{printf "\tssl-default-bind-options no-sslv3"}}
{{printf "\n"}}
{{"defaults"}}
{{printf "\tlog \tglobal"}}
{{printf "\tmode \thttp"}}
{{printf "\toption  dontlognull"}}
{{printf "\ttimeout connect 5000"}}
{{printf "\ttimeout client  50000"}}
{{printf "\ttimeout server  50000"}}
{{printf "\terrorfile 400 /etc/haproxy/errors/400.http"}}
{{printf "\terrorfile 403 /etc/haproxy/errors/403.http"}}
{{printf "\terrorfile 408 /etc/haproxy/errors/408.http"}}
{{printf "\terrorfile 500 /etc/haproxy/errors/500.http"}}
{{printf "\terrorfile 502 /etc/haproxy/errors/502.http"}}
{{printf "\terrorfile 503 /etc/haproxy/errors/503.http"}}
{{printf "\terrorfile 504 /etc/haproxy/errors/504.http"}}
{{$name := .Name}}
{{$dns := .Dns}}
{{range .Hosts -}}
{{printf "frontend f_%s-%s" $name .PortSRC}}
{{if eq .Protocol "https" -}}
{{printf "\tbind *:%s ssl crt /etc/haproxy/ssl/totino.com.br.pem crt /etc/haproxy/ssl/totino2.com.br.pem" .PortSRC}}
{{printf "\tmode http"}}
{{printf "\tacl %s_acl_https hdr(host) -i %s" $name $dns}}
{{printf "\tuse_backend b_%s-%s if %s_acl_https" $name .PortSRC $name}}
{{printf "\n"}}
{{printf "backend b_%s-%s" $name .PortSRC}}
{{printf "\tmode http"}}
{{printf "\tsource 0.0.0.0 usesrc client"}}
{{printf "\tserver minion-1 minion-1.com.br:%s check ssl verify none" .PortSRC}}
{{printf "\tserver minion-2 minion-2.com.br:%s check ssl verify none" .PortSRC}}
{{printf "\tserver minion-3 minion-3.com.br:%s check ssl verify none" .PortSRC}}
{{printf "\tlog /dev/log local0 debug"}}
{{else -}}{{if eq .Protocol "http" -}}
{{printf "\tbind *:%s" .PortSRC}}
{{printf "\tmode http"}}
{{printf "\tacl %s_acl_http hdr(host) -i %s" $name $dns}}
{{printf "\tuse_backend b_%s-%s if %s_acl_http" $name .PortSRC $name}}
{{printf "\n"}}
{{printf "backend b_%s-%s" $name .PortSRC}}
{{printf "\tmode http"}}
{{printf "\tsource 0.0.0.0 usesrc client"}}
{{printf "\thttp-request set-header Host %s" $dns}}
{{printf "\tserver minion-1 minion-1.com.br:%s check" .PortSRC}}
{{printf "\tserver minion-2 minion-2.com.br:%s check" .PortSRC}}
{{printf "\tserver minion-3 minion-3.com.br:%s check" .PortSRC}}
{{printf "\tlog /dev/log local0 debug"}}
{{else -}}
{{printf "\tbind *:%s" .PortSRC}}
{{printf "\tmode tcp"}}
{{printf "\tuse_backend b_%s-%s" $name .PortSRC}}
{{printf "\n"}}
{{printf "backend b_%s-%s" $name .PortSRC}}
{{printf "\tmode tcp"}}
{{printf "\tsource 0.0.0.0 usesrc client"}}
{{printf "\tserver minion-1 minion-1.com.br:%s check" .PortSRC}}
{{printf "\tserver minion-2 minion-2.com.br:%s check" .PortSRC}}
{{printf "\tserver minion-3 minion-3.com.br:%s check" .PortSRC}}
{{end}}
{{end}}
{{end}}`
)
