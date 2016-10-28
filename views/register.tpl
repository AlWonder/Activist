{{if .flash.error}}
	<h3>{{.flash.error}}</h3>
{{end}}
{{if .Errors}}
	{{range $rec := .Errors}}
	<h3>{{$rec}}</h3>
	{{end}}
{{end}}
<h1>Регистрация</h1>

<form method="POST">
	<p>Имя пользователя</p>
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
	<input type="submit" value="Регистрация">
</form>