<section class="participants">
<h2>Участники</h2>
	{{range $val := .Participants}}
	<article class="participant">
		<p>{{$val.FirstName}} {{$val.SecondName}} {{$val.LastName}}, {{$val.Email}}</p>
		</article>
	{{end}}
</section>