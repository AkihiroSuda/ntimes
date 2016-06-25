package main

import (
	"encoding/json"
	"io"
	"text/template"
)

type formatter struct {
	Format string
	Writer io.Writer
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

func (f *formatter) Execute(stat *Stat) error {
	return f.tmpl.Execute(f.Writer, stat)
}
