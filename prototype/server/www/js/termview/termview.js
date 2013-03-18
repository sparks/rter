angular.module('termview', ['ngResource', 'items', 'ui.bootstrap.dialog'])

.controller('TermViewCtrl', function($scope, updateItemDialog) {
	$scope.updateItemDialog = function(item){
		updateItemDialog.open(item);
	};

	$scope.mapCenter = new google.maps.LatLng(45.50745, -73.5793);

	$scope.mapOptions = {
		center: $scope.mapCenter,
		zoom: 10,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	$scope.centerAt = function(location) {
		var latlng = new google.maps.LatLng(location.coords.latitude, location.coords.longitude);
		$scope.map.setCenter(latlng);
		$scope.mapCenter = latlng;
	};
})

.directive('termview', function(Item) {
	return {
		restrict: 'E',
		scope: {
			term: "@"
		},
		templateUrl: '/template/termview/termview.html',
		controller: 'TermViewCtrl',
		link: function(scope, element, attrs) {
			if(attrs.term === undefined) {
				scope.items = Item.query();
			} else {
				scope.items = Item.query({term: attrs.term});
			}

			navigator.geolocation.getCurrentPosition(scope.centerAt);
		}
	};
});

