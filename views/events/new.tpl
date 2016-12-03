<!DOCTYPE html>

<section class="event-new">
	<h1>Новое мероприятие</h1>

	<form method="POST">
		<label for="event-name">Название</label>
		<input type="text" id="event-name" name="event-name" required><br>
		<label for="description">Описание</label><br>
		<textarea id="description" name="description" rows="10" cols="50"></textarea><br>
		<label for="event-date">Дата проведения мероприятия</label>
		<input type="date" id="event-date" name="event-date" required><br>
		<label for="event-time">Время начала</label>
		<input type="time" id="event-time" name="event-time"><br>
		<input type="submit" value="Добавить">
	</form>
</section>