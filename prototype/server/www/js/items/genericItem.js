angular.module('genericItem', ['ng', 'ui', 'taxonomy'])

.controller('FormGenericItemCtrl', function($scope, Taxonomy) {
	if($scope.item.Author === undefined) {
		$scope.item.Author = "anonymous"; //TODO: Replace with login
	}

	//This is kinda terrible
	if($scope.item.Terms !== undefined) {
		var concat = "";
		for(var i = 0;i < $scope.item.Terms.length;i++) {
			concat += $scope.item.Terms[i].Term+",";
		}
		$scope.item.Terms = concat.substring(0, concat.length-1);
	}

	$scope.tagConfig = {
		data: Taxonomy.query(),
		multiple: true,
		id: function(item) {
			return item.Term;
		},
		formatResult: function(item) {
            return item.Term;
        },
        formatSelection: function(item) {
			return item.Term;
        },
        createSearchChoice: function(term) {
			return {Term: term};
        },
        matcher: function(term, text, option) {
			return option.Term.toUpperCase().indexOf(term.toUpperCase())>=0;
        },
        initSelection: function (element, callback) {
			var data = [];
			$(element.val().split(",")).each(function () {
				data.push({Term: this});
			});
			callback(data);
		}
	};

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

.directive('formGenericItem', function(Taxonomy, $timeout) {
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

			$timeout(
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

.directive('tileGenericItem', function(Taxonomy) {
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

.directive('closeupGenericItem', function(Taxonomy) {
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