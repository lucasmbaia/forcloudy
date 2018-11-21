package services

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/lucasmbaia/forcloudy/api/datamodels"
	"github.com/lucasmbaia/forcloudy/api/models"
	"github.com/lucasmbaia/forcloudy/api/repository"
	"github.com/satori/go.uuid"
	"reflect"
	"regexp"
	"strings"
)

const (
	EMPTY_STR = ""
)

type ResourceService interface {
	Post(ctx iris.Context) (datamodels.Response, error)
	Get(ctx iris.Context) (interface{}, error)
	GetById(ctx iris.Context, id string) (interface{}, error)
	/*Print()
	  Set(params map[string]interface{}) error
	  GetFields() interface{}
	  Post() error
	  Get() (interface{}, error)*/
}

type resourceService struct {
	fields     func() interface{}
	model      func(repository.Repositorier) models.Models
	repository repository.Repositorier
}

func (r *resourceService) Post(ctx iris.Context) (response datamodels.Response, err error) {
	var (
		model  = r.model(r.repository)
		fields = r.fields()
		id     uuid.UUID
	)

	if ctx.Request().ContentLength != 0 {
		if err = ctx.ReadJSON(fields); err != nil {
			return response, err
		}
	}

	if id, err = uuid.NewV4(); err != nil {
		return response, err
	}

	if err = r.setParams(ctx.GetCurrentRoute().Path(), ctx.Params(), fields, id.String()); err != nil {
		return response, err
	}

	r.setDependences(fields)
	if err = model.Post(fields); err != nil {
		return response, err
	}

	response.ID = id.String()
	return response, nil
}

func (r *resourceService) Get(ctx iris.Context) (i interface{}, err error) {
	var (
		model  = r.model(r.repository)
		fields = r.fields()
	)

	if err = r.setParams(ctx.GetCurrentRoute().Path(), ctx.Params(), fields, EMPTY_STR); err != nil {
		return i, err
	}

	fmt.Println(fields)
	return model.Get(fields)
}

func (r *resourceService) GetById(ctx iris.Context, id string) (i interface{}, err error) {
	var (
		model  = r.model(r.repository)
		fields = r.fields()
	)

	if err = r.setParams(ctx.GetCurrentRoute().Path(), ctx.Params(), fields, id); err != nil {
		return i, err
	}

	return model.Get(fields)
}

func (r *resourceService) setParams(url string, rp *context.RequestParams, m interface{}, id string) error {
	var (
		rgx     = regexp.MustCompile(`{[^}]*}`)
		matches []string
		params  = make(map[string]interface{})
		value   string
		err     error
		body    []byte
	)

	matches = rgx.FindAllString(url, -1)
	for _, v := range matches {
		value = strings.Replace(strings.Split(v, ":")[0], "{", "", -1)
		params[value] = rp.Get(value)
	}

	if id != EMPTY_STR {
		params["ID"] = id
	}

	if body, err = json.Marshal(params); err != nil {
		return err
	}

	if err = json.Unmarshal(body, m); err != nil {
		return err
	}

	return nil
}

func (r *resourceService) setDependences(m interface{}) {
	var (
		v reflect.Value
		t reflect.Type
	)

	v = reflect.Indirect(reflect.ValueOf(m))
	t = v.Type()

	//fmt.Println(v.FieldByName("ID").String())

	for i := 0; i < v.NumField(); i++ {
		if tag, ok := t.FieldByName(t.Field(i).Name); ok {
			if tag.Tag.Get("gorm") != EMPTY_STR {
				var (
					af          string
					association []string
					fk          string
				)

				af, _ = tag.Tag.Lookup("gorm")
				association = strings.Split(af, ";")
				fk = strings.Split(association[1], ":")[1]

				if v.Field(i).Kind() == reflect.Slice {
					setArgs(v.Field(i), fk, v.FieldByName(strings.Split(association[0], ":")[1]).Interface())
				}
			}
		}
	}
}

func setArgs(v reflect.Value, fk string, values interface{}) {
	for i := 0; i < v.Len(); i++ {
		v.Index(i).FieldByName(fk).Set(reflect.ValueOf(values))
	}
}

/*func (r *resourceService) Print() {
  fmt.Println("MODEL")
  fmt.Println(r.fields)
}

func (r *resourceService) GetFields() interface{} {
  return r.fields
}

func (r *resourceService) Set(params map[string]interface{}) error {
  var (
    body []byte
    err  error
  )

  if body, err = json.Marshal(params); err != nil {
    return err
  }

  if err = json.Unmarshal(body, r.fields); err != nil {
    return err
  }

  return nil
}

func (r resourceService) Get() (interface{}, error) {
  return r.model.Get(r.fields)
}

func (r *resourceService) Post() error {
  return r.model.Post(r.fields)
}*/
