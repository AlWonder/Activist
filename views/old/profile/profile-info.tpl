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