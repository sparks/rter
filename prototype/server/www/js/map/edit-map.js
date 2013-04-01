angular.module('edit-map', [
	'ui', //Map
	'ng'  //$timeout
])

.controller('EditMapCtrl', function($scope, $timeout) {
	$scope.mapOptions = {
		center: new google.maps.LatLng(45.50745, -73.5793),
		zoom: 10,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	$scope.setCenter = function(latLng) {
		$timeout(function() {
			$scope.map.setCenter(latLng);
		}, 0);
	};

	$scope.setMarker = function(latLng) {
		if($scope.marker !== undefined) {
			$scope.marker.setPosition(latLng);
		} else {
			$scope.marker = new google.maps.Marker({
				map: $scope.map,
				position: latLng
			});
		}
	};

	$scope.mapClick = function($event) {
		$scope.item.Lat = $event.latLng.lat();
		$scope.item.Lng = $event.latLng.lng();

		$scope.setMarker($event.latLng);
	};

	$scope.$watch('[item.Lat, item.Lng]', function() {
		if($scope.item.Lat !== undefined && $scope.item.Lng !== undefined && ($scope.item.Lat !== 0 || $scope.item.Lng !== 0)) {
			$scope.item.HasGeo = true;
		} else {
			$scope.item.HasGeo = false;
		}
	}, true);
})

.directive('editMap', function($timeout) {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/map/edit-map.html',
		controller: 'EditMapCtrl',
		link: function(scope, element, attrs) {
			if(scope.item.Lat !== undefined && scope.item.Lng !== undefined) {
				var latLng = new google.maps.LatLng(scope.item.Lat, scope.item.Lng);
				scope.setCenter(latLng);
				scope.setMarker(latLng);
			} else {
				navigator.geolocation.getCurrentPosition(function(location) {
					scope.setCenter(new google.maps.LatLng(location.coords.latitude, location.coords.longitude));
				});
			}

			$timeout( //FIXME: Another map hack to render hidden maps
				function() {
					google.maps.event.trigger(scope.map, "resize");
					scope.map.setCenter(scope.map.getCenter());
				},
				0
			);
		}
	};
});
