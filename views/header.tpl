<!DOCTYPE html>

<html>
<head>
	<title>{{.Title}}</title>
	<link href="/css/bootstrap.min.css" rel="stylesheet">
	<link href="/css/style.css" rel="stylesheet">
</head>
<body>
	<header class="navbar navbar-inverse navbar-fixed-top main-header">
		<div class='container'>
			<div class='navbar-header'>
			<a href="http://localhost:8080/home" class='navbar-brand'>Активист</a>
			<button type="button" class='navbar-toggle' data-toggle='collapse' data-target='.navbar-collapse'>
				<span class='sr-only'>Открыть меню</span>
				<i class='glyphicon glyphicon-align-justify'></i>
			</button>
			</div>
			<ul class="nav navbar-nav navbar-right collapse navbar-collapse userpanel">{{if .InSession}}
			    <!--<li><p>{{.FirstName}} {{.LastName}}</p></li>-->
			    <li><a href="http://localhost:8080/profile">Профиль</a></li>
			    <li><a href="http://localhost:8080/logout">Выход</a></li>
			    {{else}}
			    <li><a href="http://localhost:8080/register">Регистрация</a></li>
			    <li><a href="http://localhost:8080/login/home">Вход</a></li>
			    {{end}}
			</ul>
		</div>
	</header>
	<main>
	<div class="container">