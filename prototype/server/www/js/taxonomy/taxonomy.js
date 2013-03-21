angular.module('taxonomy', [
	'ngResource' //$resource for taxonomoy
])

.factory('Taxonomy', function ($resource) {
	var Taxonomy = $resource(
		'/1.0/taxonomy/:Term',
		{},
		{
			update: { method: 'PUT', params:{ Term: '@Term' } }
		}
	);

	return Taxonomy;
});