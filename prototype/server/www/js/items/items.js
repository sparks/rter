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

.controller('CreateItemCtrl', function($scope, Alerter, Item) {
	$scope.item = {Type: ""};

	$scope.debug = function() {
		console.log($scope.item);
	};

	$scope.createItem = function() {
		if($scope.item.StartTime !== undefined) $scope.item.StartTime = new Date($scope.item.StartTime);
		if($scope.item.StopTime !== undefined) $scope.item.StopTime = new Date($scope.item.StopTime);

		Item.save(
			$scope.item,
			function() {
				Alerter.success("Item Created", 2000);
				$scope.item = {Type: ""};
			},
			function(e) {
				Alerter.error("There was a problem creating the item. "+"Status:"+e.status+". Reply Body:"+e.data);
				console.log(e);
			}
		);
	};
})

.directive('createItem', function(Item) {
	return {
		restrict: 'E',
		scope: {},
		templateUrl: '/template/items/create-item.html',
		controller: 'CreateItemCtrl',
		link: function(scope, element, attrs) {

		}
	};
})

.controller('UpdateItemDialogCtrl', function($scope, item, dialog) {
	$scope.dialog = dialog;
	$scope.item = item;
})

.controller('UpdateItemCtrl', function($scope, Alerter, Item) {
	$scope.debug = function() {
		console.log($scope.item);
	};

	$scope.updateItem = function() {
		if($scope.item.StartTime !== undefined) $scope.item.StartTime = new Date($scope.item.StartTime);
		if($scope.item.StopTime !== undefined) $scope.item.StopTime = new Date($scope.item.StopTime);

		Item.update(
			$scope.item,
			function() {
				Alerter.success("Item Updated", 2000);
				if($scope.dialog !== undefined) $scope.dialog.close();
			},
			function(e) {
				Alerter.error("There was a problem updating the item. "+"Status:"+e.status+". Reply Body:"+e.data);
				console.log(e);
			}
		);
	};

	$scope.cancel = function() {
		if($scope.dialog !== undefined) {
			$scope.dialog.close();
		}
	};
})

.directive('updateItem', function(Item) {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			dialog: "="
		},
		templateUrl: '/template/items/update-item.html',
		controller: 'UpdateItemCtrl',
		link: function(scope, element, attrs) {

		}
	};
});
