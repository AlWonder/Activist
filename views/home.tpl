<section class="events-home">
{{range $val := .Events}}
	<article class="event">
		<h3 class="title"><a href="/events/view/{{$val.Id}}">{{$val.Name}}</a></h3>
		<p class="description">Описание</p>
		</article>
	{{end}}
</section>
