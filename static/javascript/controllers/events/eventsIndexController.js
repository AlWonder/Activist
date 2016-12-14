angular.module('Activist')
.controller('EventsIndexController', function(Event, $scope) {
  $scope.events = Event.query();
});
