angular.module('Activist').controller('EventsShowController', function(Event, $scope, $routeParams) {
  $scope.event = Event.get({ id: $routeParams.id });
});
