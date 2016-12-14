angular.module('Activist')
.config(function($routeProvider) {
  $routeProvider

    .when('/', {
      redirectTo: '/events'
    })

    .when('/events', {
      templateUrl: '/tpl/events/index.html',
      controller: 'EventsIndexController'
    })

    .when('/events/:id', {
      templateUrl: '/tpl/events/show.html',
      controller: 'EventsShowController'
    });
});
