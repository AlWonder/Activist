$(function() { //JQuery $(document).ready() short form
	$('#tag').on('keyup', function(event) {
		event.preventDefault();
		$.ajax('//localhost:8080/tags/find', {
			type: 'POST',
			data: {
				'tag': $('#tag').val()
			},
			success: function(result) {
				$('.tags').find('.found').remove();
				$('.tags').append(result).slideDown();
			},
			timeout: 10000
		});
	});

	$('.login-form').on('submit', function(event) {
		event.preventDefault();
		$.ajax('http://127.0.0.1/login/home', {
			type: 'POST',
			data: $(this).serialize()
		});
	});
});