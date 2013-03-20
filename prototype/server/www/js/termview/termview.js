angular.module('termview', ['ngResource', 'items', 'ui.bootstrap.dialog'])

.controller('TermViewCtrl', function($scope, updateItemDialog, closeupItemDialog, ItemCache) {
	$scope.mapResized = false;

	$scope.resizeMap = function() {
		if(!$scope.mapResized) {
			$scope.mapResized = true;
			google.maps.event.trigger($scope.map, "resize");
			$scope.map.setCenter($scope.mapCenter);
		}
	};

	$scope.updateItemDialog = function(item){
		updateItemDialog.open(item).then(function() {
			$scope.updateMarkers();
		});
	};

	$scope.closeupItemDialog = function(item){
		closeupItemDialog.open(item);
	};

	$scope.mapCenter = new google.maps.LatLng(45.50745, -73.5793);

	$scope.mapOptions = {
		center: $scope.mapCenter,
		zoom: 10,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	$scope.markers = [];

	$scope.updateMarkers = function() {
		angular.forEach($scope.markers, function(v) {
			v.setMap(null);
		});

		$scope.markers = [];

		angular.forEach($scope.items, function(v) {
			if(v.Lat === undefined || v.Lng === undefined || (v.Lat === 0 && v.Lng === 0)) return;

			var m = new google.maps.Marker({
				map: $scope.map,
				position: new google.maps.LatLng(v.Lat, v.Lng)
			});

			$scope.markers.push(m);
		});
	};

	$scope.centerAt = function(location) {
		var latlng = new google.maps.LatLng(location.coords.latitude, location.coords.longitude);
		$scope.map.setCenter(latlng);
		$scope.mapCenter = latlng;
	};

	$scope.$watch('items', function(v) {
		$scope.updateMarkers();
	});
})

.directive('termview', function(ItemCache) {
	return {
		restrict: 'E',
		scope: {
			term: "="
		},
		templateUrl: '/template/termview/termview.html',
		controller: 'TermViewCtrl',
		link: function(scope, element, attrs) {
			scope.items = ItemCache.items;

			navigator.geolocation.getCurrentPosition(scope.centerAt);
		}
	};
});

