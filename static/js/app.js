$(function() { //JQuery $(document).ready() short form
	$('#tag').on('keyup', function(event) {
		var tag = $('#tag').val().trim()
		
		event.preventDefault();
		$.ajax({
			type: 'POST',
			url: '/tags/find',
			//contentType: 'application/json',
			data: { 'tag': tag },
			dataType: 'json',
			success: function(result) {
				var list = $('.tags .found');
				//list.find('li').remove();
				//alert(result);
				/*$.each(result, function(index, found) {
					list.append('<li id="' + found.id + '">' + found.name + '</li>');
				});*/
				var foundElements = $.map(result, function(found, index) {
					var listItem = $('<li></li>');
					listItem.append('<li id="' + found.id + '">' + found.name + '</li>');
					return listItem;
				});
				list.detach()
				    .html(foundElements)
				    .appendTo('.tags');
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