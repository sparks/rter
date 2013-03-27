angular.module('youtubeItem', [
	'ng', //$timeout
	'ui'  //Map
])

.controller('FormYoutubeItemCtrl', function($scope) {
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

.directive('formYoutubeItem', function($timeout) {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/youtube/form-youtube-item.html',
		controller: 'FormYoutubeItemCtrl',
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

.controller('TileYoutubeItemCtrl', function($scope) {
	$scope.youtubeID = $scope.item.ContentURI.match(/\/watch\?v=([0-9a-zA-Z].*)/)[1];
})

.directive('tileYoutubeItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/youtube/tile-youtube-item.html',
		controller: 'TileYoutubeItemCtrl',
		link: function(scope, element, attr) {

		}
	};
})

.controller('CloseupYoutubeItemCtrl', function($scope) {
	$scope.youtubeID = $scope.item.ContentURI.match(/\/watch\?v=([0-9a-zA-Z].*)/)[1];

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

.directive('closeupYoutubeItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/youtube/closeup-youtube-item.html',
		controller: 'CloseupYoutubeItemCtrl',
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