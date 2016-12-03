<!DOCTYPE html>

<section class="event-edit">
	<h1>Редактировать мероприятие</h1>

	<form method="POST">
		<label for="event-name">Название</label>
		<input type="text" id="event-name" name="event-name" value="{{.Name}}" required><br>
		<label for="description">Описание</label><br>
		<textarea id="description" name="description" rows="10" cols="50">{{.Description}}</textarea><br>
		<label for="event-date">Дата проведения мероприятия</label>
		<input type="date" id="event-date" name="event-date"  value="{{.EventDate}}"required><br>
		<label for="event-time">Время начала</label>
		<input type="time" id="event-time" name="event-time" value="{{.EventTime}}"><br>
		<input type="submit" value="Добавить">
	</form>
</section>