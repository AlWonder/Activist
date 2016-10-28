<!DOCTYPE html>

<html>
<head>
	<title>{{.Title}}</title>
</head>
<body>
<h1>Header</h1>
<a href="http://localhost:8080/home">Home</a>
<div>{{if .InSession}}
    {{.Email}} [<a href="http://localhost:8080/logout">Logout</a>|<a href="http://localhost:8080/profile">Profile</a>]
    {{else}}
    [<a href="http://localhost:8080/login/home">Login</a>]
    {{end}}
</div>