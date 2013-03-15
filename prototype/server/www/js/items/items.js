angular.module('items', ['ngResource', 'ui', 'alerts', 'genericItem', 'rawItem', 'twitterItem'])

.factory('Item', function ($resource) {
	var Item = $resource(
		'/1.0/items/:ID',
		{},
		{
			update: { method: 'PUT', params:{ ID: '@ID' } }
		}
	);

	return Item;
})

.controller('SubmitItemCtrl', function($scope, $rootScope, Alerter, Item) {
	$scope.item = {
		Type: "generic"
	};

	$scope.pushItem = function() {
		Item.save($scope.item,
			function() {
				Alerter.success("Item Created", 2000);
			},
			function(e) {
				Alerter.error("There was a problem creating the item. "+"Status:"+e.status+". Reply Body:"+e.data);
				console.log(e);
			}
		);

		$scope.newItem = {Type:"", AuthorID:null};
	};

	$scope.bang = function() {
		console.log("Hi");
		console.log($scope.createItemForm);
	};
})

.directive('submitItem', function(Item) {
	return {
		restrict: 'E',
		scope: {},
		templateUrl: '/template/items/submit-item.html',
		controller: 'SubmitItemCtrl',
		link: function(scope, element, attrs) {
			// $compile(templ)(scope)
		}
	};
});