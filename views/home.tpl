<section class="events-home">
{{range $val := .Events}}
<div class='col-sm-4'>
	<article class="event well well-sm">
		<h3 class="title"><a href="/events/view/{{$val.Id}}">{{$val.Name}}</a></h3>
		<p class="description">Описание</p>
		</article>
		</div>
	{{end}}
</section>
