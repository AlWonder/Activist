/*$(function() { //JQuery $(document).ready() short form
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
				//$.each(result, function(index, found) {
				//	list.append('<li id="' + found.id + '">' + found.name + '</li>');
				//});
				var foundElements = $.map(result, function(found, index) {
					var listItem = $('<li id="' + found.id + '"></li>');
					listItem.append(found.name);
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
});*/

(function() {
	var app = angular.module('activist', []);

	app.controller('AcceptedEventsController', ['$http', function($http){
        var that = this;
        that.events = [];

        $http.get('/json/events/accepted').then(function successCallback(data) {
            that.events = data;
        }, function errorCallback(response) {
    // called asynchronously if an error occurs
    // or server returns response with an error status.
  });
    }]);

})();
