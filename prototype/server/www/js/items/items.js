angular.module('items', [
	'ui.bootstrap',         //dialog
	'ngResource',           //$resource for Item
	'sockjs',               //sock for ItemCache
	'alerts',               //Alerts for item actions
	'taxonomy',             //For tag-selector
	'genericItem',          //generic item implementation
	'rawItem',              //raw item implementation
	'twitterItem',          //twitter item implementation
	'youtubeItem',          //YouTube item implementation
	'singleItem',          	//SingleTweet item implementation
	'streamingVideoV1Item', //Streaming Video V1 item implementation
	'comments',             //Comments dialog
	'cache'                 //CacheBuilder service
])

.factory('ItemResource', function ($resource) {
	var ItemResource = $resource(
		'/1.0/items/:ID',
		{ ID: '@ID' },
		{
			update: { method: 'PUT', params:{ ID: '@ID' } }
		}
	);

	return ItemResource;
})

.factory('ItemStream', function(SockJS) {
	return new SockJS('/1.0/streaming/items');
})

.factory('ItemCache', function (CacheBuilder, ItemResource, ItemStream) {
	return new CacheBuilder(
		"Item",
		ItemResource,
		ItemStream,
		function(a, b) {
			if(a.ID === undefined || b.ID === undefined) return false;
			if(a.ID == b.ID) return true;
			return false;
		}
	);
})

.filter('filterbyBounds', function() {
	return function(input, bounds) {
		out = [];
		for(var i = 0;i < input.length;i++) {
			if(input[i].Lat !== undefined && input[i].Lng !== undefined && bounds !== undefined) {
				if(input[i].Lat < Math.min(bounds.getNorthEast().lat(), bounds.getSouthWest().lat()) || input[i].Lat > Math.max(bounds.getNorthEast().lat(), bounds.getSouthWest().lat())) {
					//Outside via lat
				} else if(input[i].Lng < Math.min(bounds.getNorthEast().lng(), bounds.getSouthWest().lng()) || input[i].Lng > Math.max(bounds.getNorthEast().lng(), bounds.getSouthWest().lng())) {
					//Outside via lng
				} else {
					out.push(input[i]);
				}
			}
		}

		return out;
	};
})

.filter('orderByRanking', function() { //FIXME: this is n^2 probably not good
	return function(input, ranking) {
		if(ranking === undefined || ranking.length === 0) return input;

		var out = [];
		var stragglers = [];

		for(var i = 0;i < input.length;i++) {
			var found = false;
			for(var j = 0;j < ranking.length;j++) {
				if(input[i].ID == ranking[j] && out[j] === undefined) {
					found = true;
					out[j] = input[i];
					break;
				}
			}
			if(!found) {
				stragglers.push(input[i]);
			}
		}

		for(var i = 0;i < out.length;i++) {
			if(out[i] === undefined) {
				out.remove(i);
				i--;
			}
		}

		out.push.apply(out, stragglers);

		return out;
	};
})

.filter('filterByTerm', function() {
	return function(input, term) {
		if(term === "" || term === undefined) return input;
		var out = [];
		for(var i = 0;i < input.length;i++) {
			if(input[i].Terms !== undefined) {
				for(var j = 0;j < input[i].Terms.length;j++) {
					if(input[i].Terms[j].Term == term) out.push(input[i]);
				}
			}
		}
		return out;
	};
})

.controller('CreateItemCtrl', function($scope, Alerter, ItemCache) {
	var defaultType = "";
	$scope.item = {Type: defaultType};

	$scope.inProgress = false;

	$scope.createItem = function() {
		if($scope.item.StartTime !== undefined) $scope.item.StartTime = new Date($scope.item.StartTime);
		if($scope.item.StopTime !== undefined) $scope.item.StopTime = new Date($scope.item.StopTime);

		$scope.inProgress = true;

		ItemCache.create(
			$scope.item,
			function() {
				$scope.item = {Type: defaultType};
				$scope.inProgress = false;
			},
			function() {
				$scope.inProgress = false;
			}
		);
	};
})

.directive('createItem', function() {
	return {
		restrict: 'E',
		scope: {},
		templateUrl: '/template/items/create-item.html',
		controller: 'CreateItemCtrl',
		link: function(scope, element, attrs) {

		}
	};
})

.controller('UpdateItemCtrl', function($scope, Alerter, ItemCache) {
	$scope.inProgress = false;

	$scope.updateItem = function() {
		if($scope.itemCopy.StartTime !== undefined) $scope.itemCopy.StartTime = new Date($scope.itemCopy.StartTime);
		if($scope.itemCopy.StopTime !== undefined) $scope.itemCopy.StopTime = new Date($scope.itemCopy.StopTime);

		$scope.inProgress = true;

		ItemCache.update(
			$scope.itemCopy,
			function() {
				$scope.cancel();
				$scope.inProgress = false;
			},
			function(e) {
				if(e.status == 304) {
					$scope.cancel();
				}
				$scope.inProgress = false;
			}
		);
	};

	$scope.deleteItem = function() {
		ItemCache.remove(
			$scope.item,
			function() {
				$scope.cancel();
			}
		);
	};

	$scope.cancel = function() {
		if($scope.dialog !== undefined) {
			$scope.dialog.close();
		}
	};

	$scope.itemCopy = angular.copy($scope.item); //This must be here or we break the tag
})

.directive('updateItem', function() {
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
})

.controller('UpdateItemDialogCtrl', function($scope, item, dialog) {
	$scope.dialog = dialog;
	$scope.item = item;
})

.factory('UpdateItemDialog', function ($dialog) {
	return {
		open: function(item) {
			var d = $dialog.dialog({
				modalFade: false,
				backdrop: false,
				keyboard: true,
				backdropClick: false,
				resolve: {item: function() { return item; }},
				templateUrl: '/template/items/update-item-dialog.html',
				controller: 'UpdateItemDialogCtrl'
			});

			return d.open();
		}
	};
})

.controller('TileItemCtrl', function($scope) {

})

.directive('tileItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/tile-item.html',
		controller: 'TileItemCtrl',
		link: function(scope, element, attrs) {

		}
	};
})

.controller('CloseupItemCtrl', function($scope, ItemCache, UpdateItemDialog) {
	$scope.updateItem = function() {
		ItemCache.update(
			$scope.item,
			function() {
				$scope.cancel();
			},
			function(e) {
				if(e.status == 304) {
					$scope.cancel();
				}
			}
		);
	};

	$scope.editDialog = function() {
		$scope.cancel();
		UpdateItemDialog.open($scope.item);
	};

	$scope.deleteItem = function() {
		ItemCache.remove(
			$scope.item,
			function() {
				$scope.cancel();
			}
		);
	};

	$scope.cancel = function() {
		if($scope.dialog !== undefined) {
			$scope.dialog.close();
		}
	};
})

.directive('closeupItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			dialog: "="
		},
		templateUrl: '/template/items/closeup-item.html',
		controller: 'CloseupItemCtrl',
		link: function(scope, element, attrs) {

		}
	};
})

.controller('CloseupItemDialogCtrl', function($scope, item, dialog) {
	$scope.dialog = dialog;
	$scope.item = item;
})

.factory('CloseupItemDialog', function ($dialog) {
	return {
		open: function(item) {
			var d = $dialog.dialog({
				modalFade: false,
				backdrop: false,
				keyboard: true,
				backdropClick: false,
				resolve: {item: function() { return item; }},
				templateUrl: '/template/items/closeup-item-dialog.html',
				controller: 'CloseupItemDialogCtrl'
			});

			return d.open();
		}
	};
});
