var rterCRUD = angular.module('rterCRUD', ['ngResource', 'ui.bootstrap']);

rterCRUD.factory('Item', function ($resource) {
	var Item = $resource(
		'/1.0/items/:ID',
		{},
		{
			update: { method: 'PUT', params:{ ID: '@ID' } }
		}
	);

	return Item;
});

rterCRUD.factory("Term", function ($resource) {
	var Term = $resource(
		'/1.0/taxonomy/:term',
		{},
		{
			update: { method: 'PUT', params:{ term: '@Term' } }
		}
	);

	return Term;
});