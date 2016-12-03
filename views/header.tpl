<!DOCTYPE html>

<html>
<head>
	<title>{{.Title}}</title>
</head>
<body>
<header class="main-header">
	<h1>Активист</h1>
	<a href="http://localhost:8080/home">Home</a>
	<div class="userpanel">{{if .InSession}}
	    {{.FirstName}} {{.LastName}} [<a href="http://localhost:8080/logout">Logout</a>|<a href="http://localhost:8080/profile">Profile</a>]
	    {{else}}
	    [<a href="http://localhost:8080/login/home" class="login-button">Login</a>]
	    {{end}}
	</div>
</header>
<main>