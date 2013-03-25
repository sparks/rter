angular.module('genericItem', [
	'ng', //$timeout
	'ui'  //Map
])

.controller('FormGenericItemCtrl', function($scope) {
	if($scope.item.Author === undefined) {
		$scope.item.Author = "anonymous"; //TODO: Replace with login
	}

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

	$scope.setMarker = function($event) {
		if($scope.marker === undefined) {
			$scope.marker = new google.maps.Marker({
				map: $scope.map,
				position: $event.latLng
			});
		} else {
			$scope.marker.setPosition($event.latLng);
		}

		$scope.item.Lat = $event.latLng.lat();
		$scope.item.Lng = $event.latLng.lng();
	};
})

.directive('formGenericItem', function($timeout) {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/generic/form-generic-item.html',
		controller: 'FormGenericItemCtrl',
		link: function(scope, element, attr) {
			if(scope.item.Lat !== undefined && scope.item.Lng !== undefined) {
				var latLng = new google.maps.LatLng(scope.item.Lat, scope.item.Lng);
				scope.marker = new google.maps.Marker({
					map: scope.map,
					position: latLng
				});
				scope.mapCenter = latLng;
			} else {
				navigator.geolocation.getCurrentPosition(scope.centerAt);
			}

			$timeout( //FIXME: Another map hack to render hidden maps
				function() {
					google.maps.event.trigger(scope.map, "resize");
					scope.map.setCenter(scope.mapCenter);
				},
				5
			);
		}
	};
})

.controller('TileGenericItemCtrl', function($scope) {

})

.directive('tileGenericItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/generic/tile-generic-item.html',
		controller: 'TileGenericItemCtrl',
		link: function(scope, element, attr) {

		}
	};
})

.controller('CloseupGenericItemCtrl', function($scope) {
	if($scope.item.Lat !== undefined && $scope.item.Lng !== undefined) {
		$scope.mapCenter = new google.maps.LatLng($scope.item.Lat, $scope.item.Lng);
	} else {
		$scope.mapCenter = new google.maps.LatLng(45.50745, -73.5793);
	}

	$scope.mapOptions = {
		center: $scope.mapCenter,
		zoom: 10,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};
})

.directive('closeupGenericItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/generic/closeup-generic-item.html',
		controller: 'CloseupGenericItemCtrl',
		link: function(scope, element, attr) {
			if(scope.item.Lat !== undefined && scope.item.Lng !== undefined) {
				scope.marker = new google.maps.Marker({
					map: scope.map,
					position: new google.maps.LatLng(scope.item.Lat, scope.item.Lng)
				});
			}
		}
	};
});