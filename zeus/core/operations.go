package core

import (
)

type Deploy struct {
  Customer	  string
  ApplicationName string
  ImageVersion	  string
}

func DeployAppication(d Deploy) {
  var (
    image string
  )

  image = fmt.Sprintf("%s_app-%s/image:%s", d.Customer, d.ApplicationName, d.ImageVersion)
}
