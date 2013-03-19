angular.module('termview', ['ngResource', 'items', 'ui.bootstrap.dialog'])

.controller('TermViewCtrl', function($scope, updateItemDialog, Item) {
	$scope.updateItemDialog = function(item){
		updateItemDialog.open(item).then(function() {
			$scope.items = Item.query(function() { //FIXME: This causes the page to snap up as everything is rebuilt
				$scope.updateMarkers();
			});
		});
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
				scope.items = Item.query(function() {
					scope.updateMarkers();
				});
			} else {
				scope.items = Item.query({term: attrs.term});
			}

			navigator.geolocation.getCurrentPosition(scope.centerAt);
		}
	};
});

