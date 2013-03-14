angular.module('rter', ['rterCRUD', 'alerts'])

.controller('TabsCtrl', function($scope) {
	$scope.termviews = [
		{ term:"a" },
		{ term:"b" }
	];
})

.controller('RterCtrl', function($scope, Item, Alerter) {
	$scope.items = Item.query();

	$scope.addAlert = function() {
		Alerter.error("ahhh", 1000);
		Alerter.success("ahhh", 2000);
		Alerter.warn("ahhh", 3000);
		Alerter.alert({msg: "fucl"}, 4000);
		Alerter.alert({msg: "fucl"}, 5000);
	};

	$scope.pushItem = function() {
		Item.save($scope.newItem,
			function(builtItem) {
				$scope.items.push(builtItem);
			},
			function(e) {
				console.log(e);
			}
		);

		$scope.newItem = {Type:"", AuthorID:null};
	};

	$scope.getItem = function() {
		Item.get($scope.newItem,
			function(gotItem) {
				$scope.theItem = gotItem;
			},
			function(e) {
				console.log(e);
			}

		);
		$scope.newItem = {ID:null};
	};

	$scope.setUpdateItem = function(item) {
		$scope.updateItem = {ID:item.ID, AuthorID:item.AuthorID, Type: item.Type};
	};

	$scope.putItem = function() {
		Item.update($scope.updateItem,
			function(updatedItem) {
				var index = 0;
				angular.forEach(
					$scope.items,
					function(value, key){
						if(value.ID == updatedItem.ID) {
							index = key;
						}
					}
				);
				$scope.items[index] = updatedItem;
			},
			function(e) {
				console.log(e);
			}
		);

		$scope.updateItem = {Type:"", ID:null, AuthorID:null};
	};

	$scope.deleteItem = function(item) {
		Item.remove({ID: item.ID},
			function() {
				var index = 0;
				angular.forEach(
					$scope.items,
					function(value, key){
						if(value.ID == item.ID) {
							index = key;
						}
					}
				);
				$scope.items.remove(index);
			},
			function(e) {
				console.log(e);
			}
		);
	};
})

.controller('TermViewCtrl', function($scope, Item) {
	$scope.bang = function() {
		alert("Got "+$scope.term);
	};
})

.directive('termview', function(Item) {
	return {
		restrict: 'E',
		scope: {
			term: "@"
		},
		templateUrl: 'template/termview.html',
		controller: 'TermViewCtrl',
		link: function(scope, element, attrs) {
			if(attrs.term === undefined) {
				scope.items = Item.query();
			} else {
				scope.nothing = true;
			}
		}
	};
});