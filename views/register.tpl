<h1>Регистрация</h1>

<form method="POST">
	<p>E-Mail</p>
	<input type="text" name="email" required>
	<p>Пароль</p>
	<input type="text" name="password" required>
	<p>Повторите пароль</p>
	<input type="text" name="password2" required>
	<p>Группа</p>
	<select name="group">
		<option value="1">Участник</option>
		<option value="2">Организатор</option>
	</select>
	<p>Имя</p>
	<input type="text" name="first_name" required>
	<p>Отчество</p>
	<input type="text" name="second_name" required>
	<p>Фамилия</p>
	<input type="text" name="last_name" required>
	<p>Пол</p>
	<input name="gender" type="radio" value="1"> м 
	<input name="gender" type="radio" value="2"> ж <br>
	<input type="submit" value="Регистрация">
</form>