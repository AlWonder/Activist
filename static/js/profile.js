
(function() {
	var app = angular.module('user-profile', []);

	app.directive('profileInfo', function() {
		return {
			restrict: 'E',
			templateUrl: '/ng/info'
		};
	});
})();