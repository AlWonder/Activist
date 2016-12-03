<section class="article-full">
	<h3>{{.Event.Name}}</h3>
	<p class="description">{{.Event.Description}}</p>
	<p class="event-date">Дата проведения: {{dateformat .Event.EventDate "2006-01-02"}} {{if .Event.EventTime | iszero | not}}{{dateformat .Event.EventTime "15:04"}}{{end}}</p>

	{{if .InSession}}
		{{if eq .Group 1}}
		<div class="join-event">
			{{if .IsJoined}}
				<p>Вы уже участвуете в этом мероприятии</p>
			{{else}}
				<p>Хочу участвовать!</p>
				<a href="http://localhost:8080/events/join/{{.Event.Id}}?as=1">Как участник</a> | <a href="http://localhost:8080/events/join/{{.Event.Id}}?as=2">Как волонтёр</a>
			{{end}}
		</div>
		{{end}}
	{{end}}
</section>