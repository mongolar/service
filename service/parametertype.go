package service

import (
	"net/http"
)

func init() {
	AddParameterType("form", make(FormParameter{}))
	AddParameterType("url", URLParameter{})
	AddParameterType("json", JSONParameter{})
	AddParameterType("query", QueryParameter{})
}

var ParametersTypes map[string]ParameterType

type ParameterType interface {
	Get(interface{}, *http.Request, Parameter) error
	Set(interface{}, *http.Request, Parameter) error
}

func AddParameterType(key string, pt ParameterType) {
	ParametersTypes[key] = pt
}

type FormParameter struct{}

func (fp FormParameter) Get(val interface{}, r *http.Request, p Parameter) error {
	return nil
}

func (fp FormParameter) Set(val interface{}, r *http.Request, p Parameter) error {
	return nil
}

type URLParameter struct{}

func (up URLParameter) Get(val interface{}, r *http.Request, p Parameter) error {
	return nil
}

func (up URLParameter) Set(val interface{}, r *http.Request, p Parameter) error {
	return nil
}

type JSONParameter struct{}

func (jp JSONParameter) Get(val interface{}, r *http.Request, p Parameter) error {
	return nil
}

func (jp JSONParameter) Set(val interface{}, r *http.Request, p Parameter) error {
	return nil
}

type QueryParameter struct{}

func (qp QueryParameter) Get(val interface{}, r *http.Request, p Parameter) error {
	return nil
}

func (qp QueryParameter) Set(val interface{}, r *http.Request, p Parameter) error {
	return nil
}
