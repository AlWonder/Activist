$(function() {
	$('.login-button').on('click', function() {
		$.get('http://localhost:8080/login/home', function(response) {
			$('.userpanel').html(response).slideDown();
		})
	});
});