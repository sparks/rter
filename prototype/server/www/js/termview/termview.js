angular.module('termview', ['ngResource', 'items', 'ui.bootstrap.dialog'])

.filter('filterByTerm', function() {
	return function(input, term) {
		if(term === "" || term === undefined) return input;
		var out = [];
		for(var i = 0;i < input.length;i++) {
			if(input[i].Terms !== undefined) {
				for(var j = 0;j < input[i].Terms.length;j++) {
					if(input[i].Terms[j].Term == term) out.push(input[i]);
				}
			}
		}
		return out;
	};
})

.controller('TermViewCtrl', function($scope, updateItemDialog, closeupItemDialog, Item) {
	$scope.mapResized = false;

	$scope.resizeMap = function() {
		if(!$scope.mapResized) {
				google.maps.event.trigger($scope.map, "resize");
				$scope.mapResized = true;
				$scope.map.setCenter($scope.mapCenter);
		}
	};

	$scope.updateItemDialog = function(item){
		updateItemDialog.open(item).then(function() {
			$scope.items = Item.query(function() { //FIXME: This causes the page to snap up as everything is rebuilt
				$scope.updateMarkers();
			});
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

