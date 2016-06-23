package main

import (
	"bufio"
	"encoding/json"
	"text/template"
)

type formatter struct {
	Format string
	Writer *bufio.Writer
	tmpl   *template.Template
}

func (f *formatter) Init() error {
	funcs := template.FuncMap{
		"json": func(v interface{}) string {
			b, _ := json.Marshal(v)
			return string(b)
		},
	}
	var err error
	f.tmpl, err = template.New("").Funcs(funcs).Parse(f.Format)
	return err
}

func (f *formatter) Execute(report *Report) error {
	return f.tmpl.Execute(f.Writer, report)
}
