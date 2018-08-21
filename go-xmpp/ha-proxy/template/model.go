package template

import (
  "errors"
)

const (
  HTTP_SERVER = `
{{$interface := .Interface}}
{{printf "frontend all_sites_http"}}
{{printf "\tbind %s:80" $interface}}
{{printf "\tmode http"}}
{{printf "\tlog /dev/log local0 debug"}}
{{range $idx, $host := .Hosts}}
{{printf "\tacl %s_url hdr(host) eq %s" .Name .Dns}}
{{printf "\tuse_backend bac_%s if %s_url" .Name .Name}}
{{end}}

{{range .Hosts}}
{{printf "backend bac_%s" .Name}}
{{printf "\tmode http"}}
{{printf "\thttp-request set-header Host %s" .Dns}}
{{range $idx, $minion := .Minions -}}
{{printf "\tserver host-%d %s:80 check" $idx $minion}}
{{end -}}
{{end}}`

  HTTPS_SERVER = `
{{$ssl := .SSL}}
{{$interface := .Interface}}
{{printf "frontend all_sites_https"}}
{{printf "\tbind %s:443 %s" $interface $ssl}}
{{printf "\tmode http"}}
{{printf "\tlog /dev/log local0 debug"}}
{{range $idx, $host := .Hosts}}
{{printf "\tacl %s_url hdr(host) eq %s" .Name .Dns}}
{{printf "\tuse_backend bac_%s if %s_url" .Name .Name}}
{{end}}

{{range .Hosts}}
{{printf "backend bac_%s" .Name}}
{{printf "\tmode http"}}
{{printf "\thttp-request set-header Host %s" .Dns}}
{{range $idx, $minion := .Minions -}}
{{printf "\tserver host-%d %s:80 check" $idx $minion}}
{{end -}}
{{end}}`

  HTTP_MINION = `
{{$interface := .Interface}}
{{printf "frontend all_sites_http"}}
{{printf "\tbind %s:80" $interface}}
{{printf "\tmode http"}}
{{printf "\tlog /dev/log local0 debug"}}
{{range $idx, $host := .Hosts}}
{{printf "\tacl whitelist_%s src %s" .Name .Whitelist}}
{{printf "\tacl %s_url hdr(host) eq %s" .Name .Dns}}
{{printf "\tuse_backend bac_%s if %s_url whitelist_%s" .Name .Name .Name}}
{{end}}

{{range .Hosts}}
{{printf "backend bac_%s" .Name}}
{{printf "\tmode http"}}
{{printf "\tbalance roundrobin"}}
{{printf "\thttp-request set-header Host %s" .Dns}}
{{range $idx, $addr := .Address -}}
{{printf "\tserver application-%d %s check" $idx $addr}}
{{end -}}
{{end}}`

  HTTPS_MINION = `
{{$ssl := .SSL}}
{{$interface := .Interface}}
{{printf "frontend all_sites_https"}}
{{printf "\tbind %s:443 %s" $interface $ssl}}
{{printf "\tmode http"}}
{{printf "\tlog /dev/log local0 debug"}}
{{range $idx, $host := .Hosts}}
{{printf "\tacl whitelist_%s src %s" .Name .Whitelist}}
{{printf "\tacl %s_url hdr(host) eq %s" .Name .Dns}}
{{printf "\tuse_backend bac_%s if %s_url whitelist_%s" .Name .Name .Name}}
{{end}}

{{range .Hosts}}
{{printf "backend bac_%s" .Name}}
{{printf "\tmode http"}}
{{printf "\tbalance roundrobin"}}
{{printf "\thttp-request set-header Host %s" .Dns}}
{{range $idx, $addr := .Address -}}
{{printf "\tserver application-%d %s check ssl verify none" $idx $addr}}
{{end -}}
{{end}}`

  TCP_UDP_MINION = `
{{$name := .Name}}
{{$interface := .Interface}}
{{range .Hosts -}}
{{printf "frontend front_%s-%s" $name .PortSRC}}
{{printf "\tbind %s:%s" $interface .PortSRC}}
{{printf "\tmode tcp"}}
{{printf "\tacl whitelist_%s-%s src %s" $name .PortSRC .Whitelist}}
{{printf "\tuse_backend bac_%s-%s if whitelist_%s-%s" $name .PortSRC $name .PortSRC}}

{{printf "backend bac_%s-%s" $name .PortSRC}}
{{printf "\tmode tcp"}}
{{printf "\tbalance roundrobin"}}
{{range $idx, $addr := .Address -}}
{{printf "\tserver application-%d %s check" $idx $addr}}
{{end -}}
{{end -}}`

  TCP_UDP_SERVER = `
{{$name := .Name}}
{{$interface := .Interface}}
{{range .Hosts -}}
{{$port := .PortSRC}}
{{printf "frontend front_%s-%s" $name .PortSRC}}
{{printf "\tbind %s:%s" $interface .PortSRC}}
{{printf "\tmode tcp"}}
{{printf "\tuse_backend bac_%s-%s" $name .PortSRC}}

{{printf "backend bac_%s-%s" $name .PortSRC}}
{{printf "\tmode tcp"}}
{{printf "\tbalance roundrobin"}}
{{range $idx, $minion := .Minions -}}
{{printf "\tserver host-%d %s:%s check" $idx $minion $port}}
{{end -}}
{{end -}}`

	MINION = `
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


{{/*printf "frontend all_sites"}}
{{printf "\t bind 0.0.0.0:80"}}
{{printf "\t mode http"}}
{{printf "\tlog /dev/log local0 debug"}}
{{printf "\tacl whitelist src %s" .Whitelist}}
{{range $idx, $host := .Hosts -}}
{{printf "\t acl %s_url hdr(host) eq %s" .Name .Dns}}
{{printf "\t use_backend bac_%s if whitelist %s_url" .Name .Name}}
{{end -*/}}



{{else -}}{{if eq .Protocol "http" -}}
{{printf "\t bind 0.0.0.0:80"}}
{{printf "\t mode http"}}
{{printf "\tlog /dev/log local0 debug"}}
{{printf "\tacl whitelist src %s" .Whitelist}}
{{range $idx, $host := .Hosts -}}
{{printf "\t acl %s_url hdr(host) eq %s" .Name .Dns}}
{{printf "\t use_backend bac_%s if whitelist %s_url" .Name .Name}}
{{end -}}
{{/*printf "\tbind *:%s" .PortSRC}}
{{printf "\tmode http"}}
{{printf "\tlog /dev/log local0 debug"}}
{{printf "\tacl whitelist src %s" .Whitelist}}
{{printf "\acl %s_url hdr(host) -i %s:%s" $name $dns .PortSRC}}
{{printf "\tuse_backend b_%s-%s if whitelist %_url" $name .PortSRC $name}}
{{printf "\n"*/}}

{{/*printf "backend b_%s-%s" $name .PortSRC}}
{{printf "\tmode http"}}
{{printf "\tbalance roundrobin"}}
{{range $idx, $addr := .Address -}}
{{printf "\tserver application-%d %s check" $idx $addr}}
{{end -}}
{{printf "\n"*/}}

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
{{printf "\tacl %s_acl_http hdr(host) -i %s:%s" $name $dns .PortSRC}}
{{printf "\tuse_backend b_%s-%s if %s_acl_http" $name .PortSRC $name}}
{{printf "\n"}}
{{printf "backend b_%s-%s" $name .PortSRC}}
{{printf "\tmode http"}}
{{printf "\tsource 0.0.0.0 usesrc client"}}
{{printf "\thttp-request set-header Host %s:%s" $dns .PortSRC}}
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

func ModelConf(m string) (string, error) {
  switch m {
  case "minion-http":
    return HTTP_MINION, nil
  case "minion-https":
    return HTTPS_MINION, nil
  case "server-http":
    return HTTP_SERVER, nil
  case "server-https":
    return HTTPS_SERVER, nil
  case "minion-tcpudp":
    return TCP_UDP_MINION, nil
  case "server-tcpudp":
    return TCP_UDP_SERVER, nil
  default:
    return "", errors.New("Model reported is unknown")
  }
}
