angular.module('items', ['ngResource', 'ui', 'ui.bootstrap', 'alerts', 'genericItem', 'rawItem', 'twitterItem'])

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

.factory('ItemCache', function ($rootScope, $timeout, Item, Alerter) {
	function ItemCache() {
		var self = this;

		this.items = [];

		this.refresh = function() {
			Item.query(function(newItems) { //NOTE : Clearing everything causes ui-sortable to freakout
				for(var i = 0;i < newItems.length;i++) {
					var found = false;
					for(var j = 0;i < self.items.length;j++) {
						if(self.items[j].ID == newItems[i].ID) {
							found = true;
							break;
						}
					}
					if(!found) {
						self.items.push(newItems[i]);
					}
				}
				$timeout(self.refresh, 500);
			});
		};

		this.create = function(item, sucess, failure) {
			Item.save(
				item,
				function() {
					Alerter.success("Item Created", 2000);
					self.items.push(item);
					if(angular.isFunction(sucess)) sucess();
				},
				function(e) {
					Alerter.error("There was a problem creating the item. "+"Status:"+e.status+". Reply Body:"+e.data);
					console.log(e);
					if(angular.isFunction(failure)) failure(e);
				}
			);
		};

		this.update = function(item, sucess, failure) {
			//If the item instanct is already in the array, I assume that you will externally handle 
			//rollback if the update fails. If the item instance is not in the array, this function
			//provides rollback if the update fails
			var handleRollback = this.items.indexOf(item) == -1 ? true : false;
			var present = false;
			var oldItem;

			if(handleRollback) {
				for(var i = 0;i < this.items.length;i++) {
					if(this.items[i].ID == item.ID) {
						oldItem = this.items[i];
						this.items[i] = item;
						present = true;
						break;
					}
				}
				if(!present) { //Odd, I guess we'll add as a new item
					this.items.push(item);
				}
			}

			Item.update(
				item,
				function() {
					Alerter.success("Item Updated", 2000);
					if(angular.isFunction(sucess)) sucess();
					//TODO: Make a revert mechanism here?
				},
				function(e) {
					if(e.status == 304) {
						Alerter.warn("Nothing was changed.", 2000);
					} else {
						Alerter.error("There was a problem updating the item. "+"Status:"+e.status+". Reply Body:"+e.data);
						if(handleRollback && present) {
							for(var i = 0;i < self.items.length;i++) {
								if(self.items[i].ID == oldItem.ID) {
									self.items[i] = oldItem;
									break;
								}
							}
						}
						console.log(e);
					}
					if(angular.isFunction(failure)) failure(e);
				}
			);
		};

		this.remove = function(item, sucess, failure) {
			Item.remove(
				item,
				function() {
					Alerter.success("Item Deleted", 2000);
					var indexOfItem = self.items.indexOf(item);
					if(indexOfItem == -1) {
						for(var i = 0;i < self.items.length;i++) {
							if(self.items[i].ID == item.ID) {
								self.items.remove(i);
								break;
							}
						}
					} else {
						self.items.remove(indexOfItem);
					}
					if(angular.isFunction(sucess)) sucess();
				},
				function(e) {
					Alerter.error("There was a problem deleting the item. "+"Status:"+e.status+". Reply Body:"+e.data);
					console.log(e);
					if(angular.isFunction(failure)) failure();
				}
			);
		};

		this.refresh();

	}

	return new ItemCache();
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

	$scope.debug = function() {
		console.log($scope.item);
	};

	$scope.createItem = function() {
		if($scope.item.StartTime !== undefined) $scope.item.StartTime = new Date($scope.item.StartTime);
		if($scope.item.StopTime !== undefined) $scope.item.StopTime = new Date($scope.item.StopTime);

		ItemCache.create(
			$scope.item,
			function() {
				$scope.item = {Type: defaultType};
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

.controller('UpdateItemCtrl', function($scope, Alerter, ItemCache) {
	$scope.debug = function() {
		console.log("original item", $scope.item);
		console.log("copy item", $scope.item);
	};

	$scope.updateItem = function() {
		if($scope.itemCopy.StartTime !== undefined) $scope.itemCopy.StartTime = new Date($scope.itemCopy.StartTime);
		if($scope.itemCopy.StopTime !== undefined) $scope.itemCopy.StopTime = new Date($scope.itemCopy.StopTime);

		ItemCache.update(
			$scope.itemCopy,
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

	$scope.deleteItem = function() {
		ItemCache.remove(
			item,
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
			scope.itemCopy = angular.copy(scope.item);
		}
	};
})

.controller('UpdateItemDialogCtrl', function($scope, item, dialog) {
	$scope.dialog = dialog;
	$scope.item = item;
})

.factory('updateItemDialog', function ($dialog) {
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

.directive('tileItem', function(Item) {
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

.controller('CloseupItemCtrl', function($scope) {
	$scope.cancel = function() {
		if($scope.dialog !== undefined) {
			$scope.dialog.close();
		}
	};

	$scope.debug = function() {
		console.log($scope.item);
	};
})

.directive('closeupItem', function(Item) {
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

.factory('closeupItemDialog', function ($dialog) {
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
