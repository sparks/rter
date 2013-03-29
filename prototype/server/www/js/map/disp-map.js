angular.module('disp-map', [
	'ui', //Map
	'ng'  //$timeout
])

.controller('DispMapCtrl', function($scope, $timeout) {
	$scope.targetHeading = 0;

	$scope.mapOptions = {
		zoom: 16,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	$scope.setCenter = function(latLng) {
		$timeout(function() {
			$scope.map.setCenter(latLng);
		}, 0);
	};

	$scope.rebuildMarker = function() {
		if($scope.enableMarker === undefined) return;
		if($scope.item.Lat === undefined || $scope.item.Lng === undefined) return;

		if($scope.marker !== undefined) {
			$scope.marker.setPosition(new google.maps.LatLng($scope.item.Lat, $scope.item.Lng));
		} else {
			$scope.marker = new google.maps.Marker({
				position: new google.maps.LatLng($scope.item.Lat, $scope.item.Lng),
				map: $scope.map,
				clickable: false,
				zIndex: 2
			});
		}
	};

	$scope.rebuildFov = function() {
		if($scope.enableFov === undefined) return;
		if($scope.item.Heading === undefined) return;

		var vshape = {
			path: google.maps.SymbolPath.BACKWARD_OPEN_ARROW,
			strokeColor: 'black',
			fillOpacity: 0.2,
			fillColor: 'black',
			strokeWeight: 3,
			scale: 100,
			rotation: $scope.item.Heading
		};

		if($scope.markerFOV !== undefined) {
			$scope.markerFOV.setIcon(vshape);
			$scope.markerFOV.setPosition(new google.maps.LatLng($scope.item.Lat, $scope.item.Lng));
		} else {
			$scope.markerFOV = new google.maps.Marker({
				position: new google.maps.LatLng($scope.item.Lat, $scope.item.Lng),
				icon: vshape,
				map: $scope.map,
				clickable: false,
				zIndex: 1
			});
		}
	};

	$scope.rebuildDir = function() {
		if($scope.enableFov === undefined) return;
		if($scope.targetHeading === undefined) return;

		var arrow = {
			path: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
			strokeColor: 'black',
			fillOpacity: 1,
			fillColor: '#AAA',
			strokeWeight: 3,
			scale: 10,
			rotation: $scope.targetHeading,
			anchor: new google.maps.Point(0, 10)
		};

		if($scope.markerDir !== undefined) {
			$scope.markerDir.setIcon(arrow);
			$scope.markerDir.setPosition(new google.maps.LatLng($scope.item.Lat, $scope.item.Lng));
		} else {
			$scope.markerDir = new google.maps.Marker({
				position: new google.maps.LatLng($scope.item.Lat, $scope.item.Lng),
				icon: arrow,
				map: $scope.map,
				clickable: false,
				zIndex: 2
			});
		}
	};

	$scope.mapClick = function($event) {
		if($scope.enableDir === undefined) return;

		$scope.targetHeading = google.maps.geometry.spherical.computeHeading($scope.map.getCenter(), $event.latLng);
		$scope.rebuildDir();
	};

	$scope.$watch('[item.Lat, item.Lng]', function() {
		$scope.setCenter(new google.maps.LatLng($scope.item.Lat, $scope.item.Lng));
		$scope.rebuildMarker();
		$scope.rebuildFov();
		$scope.rebuildDir();
	}, true);

	$scope.$watch('item.Heading', function() {
		$scope.rebuildFov();
	});

})

.directive('dispMap', function() {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			enableFov: "@",
			enableDir: "@",
			enableMarker: "@"
		},
		templateUrl: '/template/map/disp-map.html',
		controller: 'DispMapCtrl',
		link: function(scope, element, attrs) {

		}
	};
});
