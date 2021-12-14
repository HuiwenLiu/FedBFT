package main

import (
	"html/template"
	"os"
)

func main() {
	tmplString := "Hello, {{.}}"
	name := "Donald"
	tmpl := template.Must(template.New("hello").Parse(tmplString))
	err := tmpl.Execute(os.Stdout, name)
	if err != nil{
		panic("Could not execute template")
	}

	tmplString = "{{ if .}}Hello{{ end }}"
	tmpl = template.Must(template.New("hello2").Parse(tmplString))
	err = tmpl.Execute(os.Stdout,0)
	_ = tmpl.Execute(os.Stdout,1)

	names := []string{"Donald", "Bob", "Brodie"}
	tmplString = "{{ range .}}Hello, {{.}}{{ end }}"
	tmpl = template.Must(template.New("range").Parse(tmplString))
	_ = tmpl.Execute(os.Stdout,names)


	//tmpl, err := template.New("hello").Parse(tmplString);
	//if err != nil{
	//	panic("Could not parse the template")
	//}
}