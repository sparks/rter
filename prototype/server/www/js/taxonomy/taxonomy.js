angular.module('taxonomy', ['ngResource'])

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