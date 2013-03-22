angular.module('items', [
	'ui.bootstrap', //dialog
	'ngResource',   //$resource for Item
	'sockjs',       //sock for ItemCache
	'alerts',       //Alerts for item actions
	'taxonomy',     //For tag form
	'genericItem',  //generic item implementation
	'rawItem',      //raw item implementation
	'twitterItem'   //twitter item implementation
])

.factory('ItemResource', function ($resource) {
	var ItemResource = $resource(
		'/1.0/items/:ID',
		{},
		{
			update: { method: 'PUT', params:{ ID: '@ID' } }
		}
	);

	return ItemResource;
})

.factory('ItemStream', function(SockJS) {
	return new SockJS('/1.0/streaming/items');
})

.factory('ItemCache', function ($rootScope, $timeout, ItemResource, ItemStream, Alerter) {
	function ItemCache() {
		var self = this;

		this.stream = ItemStream;

		this.stream.onopen = function() {
			console.log('SockJS Item Stream Open');
		};

		this.stream.onmessage = function(e) {
			var bundle = e.data;
			if(bundle.action == "create") {

			} else if(bundle.action == "update") {

			} else if(bundle.action == "delete") {

			} else {
				console.log("Malformed message in Item Stream");
				console.log(e);
			}
			// self.messages.push({text:bundle.Message, ID:bundle.ID});
			// $rootScope.$digest();
		};

		this.stream.onclose = function() {
			console.log('SockJS Item Stream Closed');
		};

		this.items = [];

		this.refresh = function() {
			ItemResource.query(function(newItems) { //NOTE : Clearing everything causes ui-sortable to freakout
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
				// $timeout(self.refresh, 500);
			});
		};

		this.refresh();

		this.create = function(item, sucess, failure) {
			ItemResource.save(
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

.controller('CreateItemCtrl', function($scope, Alerter, ItemCache, Taxonomy) {
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

	//This is kinda terrible
	if($scope.item.Terms !== undefined) {
		var concat = "";
		for(var i = 0;i < $scope.item.Terms.length;i++) {
			concat += $scope.item.Terms[i].Term+",";
		}
		$scope.item.Terms = concat.substring(0, concat.length-1);
	}

	$scope.tagConfig = {
		data: Taxonomy.query(),
		multiple: true,
		id: function(item) {
			return item.Term;
		},
		formatResult: function(item) {
            return item.Term;
        },
        formatSelection: function(item) {
			return item.Term;
        },
        createSearchChoice: function(term) {
			return {Term: term};
        },
        matcher: function(term, text, option) {
			return option.Term.toUpperCase().indexOf(term.toUpperCase())>=0;
        },
        initSelection: function (element, callback) {
			console.log(element);
			var data = [];
			$(element.val().split(",")).each(function () {
				data.push({Term: this});
			});
			callback(data);
		}
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

.controller('UpdateItemCtrl', function($scope, Alerter, ItemCache, Taxonomy) {
	$scope.debug = function() {
		console.log("original item", $scope.item);
		console.log("copy item", $scope.itemCopy);
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

	//This is kinda terrible
	$scope.itemCopy = angular.copy($scope.item);

	if($scope.itemCopy.Terms !== undefined) {
		var concat = "";
		for(var i = 0;i < $scope.itemCopy.Terms.length;i++) {
			concat += $scope.itemCopy.Terms[i].Term+",";
		}
		$scope.itemCopy.Terms = concat.substring(0, concat.length-1);
	}

	$scope.tagConfig = {
		data: Taxonomy.query(),
		multiple: true,
		id: function(item) {
			return item.Term;
		},
		formatResult: function(item) {
			return item.Term;
		},
		formatSelection: function(item) {
			return item.Term;
		},
		createSearchChoice: function(term) {
			return {Term: term};
		},
		matcher: function(term, text, option) {
			return option.Term.toUpperCase().indexOf(term.toUpperCase())>=0;
		},
		initSelection: function (element, callback) {
			var data = [];
			$(element.val().split(",")).each(function (v, a) {
				data.push({Term: a});
			});
			console.log(data);
			callback(data);
		}
	};
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
