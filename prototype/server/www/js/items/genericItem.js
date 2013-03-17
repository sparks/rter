angular.module('genericItem', ['ng', 'taxonomy'])

.controller('SubmitGenericItemCtrl', function($scope, Taxonomy) {
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
})

.directive('submitGenericItem', function(Taxonomy) {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/submit-generic-item.html',
		controller: 'SubmitGenericItemCtrl',
		link: function(scope, element, attr) {
			navigator.geolocation.getCurrentPosition(scope.centerAt);

			if(scope.item.Lat !== undefined && scope.item.Lng !== undefined) {
				scope.marker = new google.maps.Marker({
					map: scope.myMap,
					position: new google.maps.LatLng(scope.item.Lat, scope.item.Lng)
				});
			}
		}
	};
});
