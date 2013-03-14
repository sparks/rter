angular.module('items', ['ngResource']).

factory('Item', function ($resource) {
	var Item = $resource(
		'/1.0/items/:ID',
		{},
		{
			update: { method: 'PUT', params:{ID:'@ID'} }
		}
	);

	return Item;
});
