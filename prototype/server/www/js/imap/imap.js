angular.module('imap', [
	'ui' //Map
])

.controller('ImapCtrl', function($scope) {
	$scope.mapCenter = new google.maps.LatLng(45.50745, -73.5793);

	$scope.mapOptions = {
		center: $scope.mapCenter,
		zoom: 16,
		mapTypeId: google.maps.MapTypeId.ROADMAP,
		draggable: false,
		zoomControl: true,
		scrollwheel: false
	};

	$scope.centerAt = function(location) {
		var latLng = new google.maps.LatLng(location.coords.latitude, location.coords.longitude);
		$scope.map.setCenter(latLng);
		$scope.mapCenter = latLng;
	};

	$scope.$watch('item.Heading', function() {
		if($scope.item.Heading === undefined) return;

		if($scope.FOV !== undefined) $scope.FOV.setMap(null);

		var vshape = {
			path: google.maps.SymbolPath.BACKWARD_OPEN_ARROW,
			strokeColor: 'black',
			fillOpacity: 0.2,
			fillColor: 'black',
			strokeWeight: 3,
			scale: 100,
			rotation: $scope.item.Heading
		};

		$scope.FOV = new google.maps.Marker({
			position: $scope.map.getCenter(),
			icon: vshape,
			map: $scope.map,
			clickable: false,
			zIndex: 1
		});
	});

	$scope.mapClick = function(e) {
		if($scope.pointer !== undefined) $scope.pointer.setMap(null);

		var arrow = {
			path: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
			strokeColor: 'black',
			fillOpacity: 1,
			fillColor: '#AAA',
			strokeWeight: 3,
			scale: 10,
			rotation: google.maps.geometry.spherical.computeHeading($scope.map.getCenter(), e.latLng),
			anchor: new google.maps.Point(0, 10)
		};

		$scope.pointer = new google.maps.Marker({
			position: $scope.map.getCenter(),
			icon: arrow,
			map: $scope.map,
			clickable: false,
			zIndex: 2
		});
	};
})

.directive('imap', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/imap/imap.html',
		controller: 'ImapCtrl',
		link: function(scope, element, attrs) {
			if(scope.item.Lat !== undefined && scope.item.Lng !== undefined) {
				var latLng = new google.maps.LatLng(scope.item.Lat, scope.item.Lng);
				scope.mapCenter = latLng;
				scope.map.setCenter(latLng);
			} else {
				navigator.geolocation.getCurrentPosition(scope.centerAt);
			}
		}
	};
});