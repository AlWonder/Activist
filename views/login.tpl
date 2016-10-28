{{if .flash.error}}
	<h3>{{.flash.error}}</h3>
{{end}}
{{if .Errors}}
	{{range $rec := .Errors}}
	<h3>{{$rec}}</h3>
	{{end}}
{{end}}
<h1>Вход</h1>

<form method="POST">
	<p>Имя пользователя</p>
	<input type="text" name="email" required>
	<p>Пароль</p>
	<input type="password" name="password" required>
	
	<input type="submit" value="Вход">
</form>