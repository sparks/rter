angular.module('genericItem', [])

.controller('SubmitGenericItemCtrl', function($scope) {
	$scope.item.AuthorID = 1; //TODO: Replace with login

	$scope.mapCenter = new google.maps.LatLng(45.50745, -73.5793);

	$scope.mapOptions = {
		center: $scope.mapCenter,
		zoom: 10,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	$scope.centerAt = function(location) {
		var latlng = new google.maps.LatLng(location.coords.latitude, location.coords.longitude);
		$scope.myMap.setCenter(latlng);
		$scope.mapCenter = latlng;
	};

	$scope.setMarker = function($event) {
		if($scope.marker === undefined) {
			$scope.marker = new google.maps.Marker({
				map: $scope.myMap,
				position: $event.latLng
			});
		} else {
			$scope.marker.setPosition($event.latLng);
		}

		$scope.item.Lat = $event.latLng.lat();
		$scope.item.Lng = $event.latLng.lng();
	};

	$scope.bang = function() {
		console.log($scope.item);
		console.log($scope.form);
	};
})

.directive('submitGenericItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		controller: 'SubmitGenericItemCtrl',
		templateUrl: '/template/items/submit-generic-item.html',
		link: function(scope, element, attr) {
			navigator.geolocation.getCurrentPosition(scope.centerAt);
		}
	};
});