angular.module('Activist').factory('Event', function($resource) {
  return $resource('/events/:id', {id: "@id"}, {

  });
});
