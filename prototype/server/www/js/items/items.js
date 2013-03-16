angular.module('items', ['ngResource', 'ui', 'alerts', 'genericItem', 'rawItem', 'twitterItem'])

.factory('Item', function ($resource) {
	var Item = $resource(
		'/1.0/items/:Username',
		{},
		{
			update: { method: 'PUT', params:{ Username: '@Username' } }
		}
	);

	return Item;
})

.controller('SubmitItemCtrl', function($scope, $rootScope, Alerter, Item) {
	var defaultType = "twitter";
	$scope.item = {
		Type: defaultType
	};

	$scope.logModel = function() {
		console.log($scope.item);
		// console.log($scope.createItemForm);
		console.log(JSON.stringify($scope.item));
	};

	$scope.pushItem = function() {
		if($scope.item.StartTime !== undefined) {
			$scope.item.StartTime = new Date($scope.item.StartTime);
		}
		if($scope.item.StopTime !== undefined) {
			$scope.item.StopTime = new Date($scope.item.StopTime);
		}
		Item.save($scope.item,
			function() {
				Alerter.success("Item Created", 2000);

				$scope.item = {
					Type: defaultType
				};
			},
			function(e) {
				Alerter.error("There was a problem creating the item. "+"Status:"+e.status+". Reply Body:"+e.data);
				console.log(e);
			}
		);
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
