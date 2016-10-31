<h2>Участники</h2>

	{{range $val := .Participants}}
		<p>{{$val.FirstName}} {{$val.SecondName}} {{$val.LastName}}, {{$val.Email}}</p>
	{{end}}