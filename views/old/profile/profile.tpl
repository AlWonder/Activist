<profile-info></profile-info>

<section class="profile">
<h1>Мой профиль</h1>
<p>E-Mail: <<<.Email>>></p>
<p><<<.FirstName>>> <<<.SecondName>>> <<<.LastName>>></p>
<p>Пол: 
<<<if eq .Gender 0>>>
Неизвестен
<<<else if eq .Gender 1>>>
Мужской
<<<else if eq .Gender 2>>>
Женский
<<<end>>></p>
<p>Группа:
<<<if eq .Group 1>>>
Участник
<<<else if eq .Group 2>>>
Организатор
<<<else if eq .Group 3>>>
Модератор
<<<else if eq .Group 4>>>
Администратор
<<<end>>></p>

<a href="http://localhost:8080/profile/changepwd">Сменить пароль</a>

</section>

<<<if eq .Group 1>>>
	<section class="participated-events" ng-controller="AcceptedEventsController as accepted">
	<h2>Где я участвую</h2>
		<article class="event-full" ng-repeat="event in events">
			<h3>{{ event.name }}</h3>
			<p class="description">{{ event.description }}</p>
			<p class="event-date">Дата проведения: {{ event.event_date }} {{ event.event_time }} </p>
			<a class="deny" href="http://localhost:8080/events/deny/{{ event.id }}">Больше не хочу участвовать</a>
		</article>
	</section>
<<<end>>>

<<<if eq .Group 2>>>
	<section class="my-events">
	<h2>Мои мероприятия</h2>

	<a type='button' class='btn btn-logout btn-default' href="/events/new"><i class="glyphicon glyphicon-pencil"></i> Новое мероприятие</a>

	<<<range $val := .Events>>>
	<article class="event-full">
		<h3><<<$val.Name>>></h3>
		<p class="description"><<<html2str $val.Description>>></p>
		<p class="event-date">Дата проведения: <<<dateformat $val.EventDate "2006-01-02">>> <<<if $val.EventTime | iszero | not>>><<<dateformat $val.EventTime "15:04">>><<<end>>></p>
		<div class="event-control">
			<p><a class="show-participants" href="http://localhost:8080/events/participants/<<<$val.Id>>>">Участники</a></p>
			<a class="event-edit" href="http://localhost:8080/events/edit/<<<$val.Id>>>">Редактировать</a> | 
			<a class="event-delete" href="http://localhost:8080/events/delete/<<<$val.Id>>>">Удалить</a>
		</div>
	</article>
	<<<end>>>
	</section>
<<<end>>>