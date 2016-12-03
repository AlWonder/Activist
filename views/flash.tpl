{{if .flash.error}}
<section class="flash error">
	<h3>{{.flash.error}}</h3>
	</section>
{{end}}
{{if .Errors}}
<section class="flash error">
	{{range $rec := .Errors}}
	<h3>{{$rec}}</h3>
	{{end}}
</section>
{{end}}