package main

const (
	taskList string = `Here are your tasks:
{{range .}}_{{.ShortID | printf "%-3d"}} {{range .Tags}}{{. | printf "+%s"}} {{end}}{{if .Project}}{{.Project | printf "prj:%s"}}{{end}}_
{{.Description}}
{{end}}`
)
