angular.module('items', [
	'ui.bootstrap', //dialog
	'ngResource',   //$resource for Item
	'sockjs',       //sock for ItemCache
	'alerts',       //Alerts for item actions
	'taxonomy',     //For tag-selector
	'genericItem',  //generic item implementation
	'rawItem',      //raw item implementation
	'twitterItem',  //twitter item implementation
	'youtubeItem',  //YouTube item implementation
	'comments'      //Comments dialog
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

.factory('ItemCache', function ($rootScope, ItemResource, ItemStream, Alerter) {
	function ItemCache() {
		var self = this;

		this.items = [];
		this.stream = ItemStream;

		function addUpdateItem(item) {
			var found = false;

			for(var i = 0;i < self.items.length;i++) {
				if(self.items[i].ID == item.ID) {
					for (var key in item) {
						self.items[i][key] = item[key];
					}
					found = true;
					break;
				}
			}

			if(!found) self.items.push(item);
		}

		function removeItem(item) {
			for(var i = 0;i < self.items.length;i++) {
				if(self.items[i].ID == item.ID) {
					self.items.remove(i);
					break;
				}
			}
		}

		this.stream.onopen = function() {

		};

		this.stream.onmessage = function(e) {
			var bundle = e.data;

			if(bundle.Action == "create" || bundle.Action == "update") {
				//Often if the user created the item, it will already be in place so treat as an update
				addUpdateItem(bundle.Item);
			} else if(bundle.Action == "delete") {
				removeItem(bundle.Item);
			} else {
				console.log("Malformed message in Item Stream");
				console.log(e);
			}

			$rootScope.$digest();
		};

		this.stream.onclose = function() {

		};

		this.init = function() {
			ItemResource.query(
				function(newItems) {
					self.items.length = 0;
					for(var i = 0;i < newItems.length;i++) {
						addUpdateItem(newItems[i]);
					}
				},
				function(e) {
					console.log("Couldn't load items");
					console.log(e);
				}
			);
		};

		this.init();

		this.create = function(item, sucess, failure) {
			ItemResource.save(
				item,
				function(data) {
					//Do not add the item here since it has no ID, it will be added by the websocket callback
					Alerter.success("Item Created", 2000);
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

			ItemResource.update(
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
			ItemResource.remove(
				item,
				function() {
					Alerter.success("Item Deleted", 2000);
					removeItem(item);
					if(angular.isFunction(sucess)) sucess();
				},
				function(e) {
					Alerter.error("There was a problem deleting the item. "+"Status:"+e.status+". Reply Body:"+e.data);
					console.log(e);
					if(angular.isFunction(failure)) failure();
				}
			);
		};
	}

	return new ItemCache();
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
