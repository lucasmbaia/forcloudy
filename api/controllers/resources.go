package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/lucasmbaia/forcloudy/api/services"
	"regexp"
	"strings"
)

type ResourceController struct {
	Services services.ResourceService
	Ctx      iris.Context
}

func (r *ResourceController) BeginRequest(ctx iris.Context) {
	fmt.Println(ctx.URLParams())
	/*  var (
	      fields  = r.Services.GetFields()
	      err	    error
	    )

	    if ctx.Request().ContentLength != 0 {
	      if err = ctx.ReadJSON(fields); err != nil {
	        fmt.Println(err)
	      }
	    }

	    if err = r.setParams(ctx.GetCurrentRoute().Path(), ctx.Params()); err != nil {
	      fmt.Println(err)
	    }

	    r.Services.Print()*/
	/*ctx.Params().Visit(func(name string, value string) {
	  fmt.Println(name, value)
	})*/
}

func (r *ResourceController) EndRequest(ctx iris.Context) {
	ctx.ContentType("application/json")
}

func (r *ResourceController) setParams(url string, ctx *context.RequestParams) error {
	var (
		rgx     = regexp.MustCompile(`{[^}]*}`)
		matches []string
		params  = make(map[string]interface{})
		value   string
		//err     error
	)

	matches = rgx.FindAllString(url, -1)
	for _, v := range matches {
		value = strings.Replace(strings.Split(v, ":")[0], "{", "", -1)
		params[value] = ctx.Get(value)
	}

	/*if err = r.Services.Set(params); err != nil {
	  return err
	}*/

	return nil
}
