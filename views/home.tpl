{{if .flash.error}}
<h3>{{.flash.error}}</h3>
{{end}}

{{range $val := .Events}}
		<h3><a href="/events/view/{{$val.Id}}">{{$val.Name}}</a></h3>
	{{end}}