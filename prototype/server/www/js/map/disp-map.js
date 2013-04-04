angular.module('disp-map', [
	'ui',  //Map
	'ng',   //$timeout
	'auth' //UserDirectionResource
])

.controller('DispMapCtrl', function($scope, $timeout, UserDirectionCache) {
	$scope.mapOptions = {
		zoom: 16,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	$scope.directionCache = new UserDirectionCache($scope.item.Author);
	$scope.userDir = $scope.directionCache.direction;

	$scope.$on("$destroy", function() {
		$scope.directionCache.close();
	});

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
		if($scope.item.Heading === undefined) {
			$scope.item.Heading = 0;
		}
		if($scope.item.Lat === undefined || $scope.item.Lng === undefined) return;

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
		if($scope.enableDir === undefined) return;
		if($scope.item.Lat === undefined || $scope.item.Lng === undefined) return;

		var arrow = {
			path: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
			strokeColor: 'black',
			fillOpacity: 1,
			fillColor: '#AAA',
			strokeWeight: 3,
			scale: 10,
			rotation: $scope.userDir.Heading,
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

		$scope.userDir.Heading = google.maps.geometry.spherical.computeHeading($scope.map.getCenter(), $event.latLng);
		$scope.rebuildDir();

		$scope.directionCache.update($scope.userDir);
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

	$scope.$watch('[enableFov, enableDir, enableMarker]', function() {
		$scope.rebuildMarker();
		$scope.rebuildFov();
		$scope.rebuildDir();
	}, true);

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
