{{if .Tags}}
<div class="found">
{{range $val := .Tags}}
<article class="tag">
	<p>{{$val.Name}}</p>
	</article>
{{end}}
</div>
{{end}}