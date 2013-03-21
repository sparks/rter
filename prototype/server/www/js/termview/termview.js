angular.module('termview', [
	'ui',          //ui-sortable and map
	'items'        //ItemCache to load items into termview, various itemDialog services
])

.factory('TermViewRemote', function () {
	function TermViewRemote() {
		this.termViews = [];

		this.addTermView = function(term) {
			for(var i = 0;i < this.termViews.length;i++) {
				if(this.termViews[i].term.Term == term.Term) {
					this.termViews[i].active = true;
					return;
				}
			}

			if(term.Term !== "") {
				this.termViews.push({term: term, heading: term.Term, active: true});
			} else {
				this.termViews.push({term: term, heading: "all", active: true});
			}
		};

		this.removeTermView = function(term) {
			for(var i = 0;i < this.termViews.length;i++) {
				if(this.termViews[i].term.Term == term.Term) {
					this.termViews.remove(i);
					return true;
				}
			}

			return false;
		};
	}

	return new TermViewRemote();
})

.controller('TermViewCtrl', function($scope, UpdateItemDialog, CloseupItemDialog, TermViewRemote) {
	$scope.mapResized = false;

	$scope.close = function() {
		TermViewRemote.removeTermView($scope.term);
	};

	$scope.resizeMap = function() { //FIXME: Another map hack to render hidden maps
		if(!$scope.mapResized) {
			$scope.mapResized = true;
			google.maps.event.trigger($scope.map, "resize");
			$scope.map.setCenter($scope.mapCenter);
		}
	};

	$scope.updateItemDialog = function(item){
		UpdateItemDialog.open(item).then(function() {
			$scope.updateMarkers();
		});
	};

	$scope.closeupItemDialog = function(item){
		CloseupItemDialog.open(item);
	};

	$scope.mapCenter = new google.maps.LatLng(45.50745, -73.5793);

	$scope.mapOptions = {
		center: $scope.mapCenter,
		zoom: 10,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	$scope.$watch('items', function(a, b) {
		$scope.updateMarkers();
	}, true);

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

.directive('termview', function(ItemCache) {
	return {
		restrict: 'E',
		scope: {
			term: "="
		},
		templateUrl: '/template/termview/termview.html',
		controller: 'TermViewCtrl',
		link: function(scope, element, attrs) {
			scope.items = ItemCache.items;
			navigator.geolocation.getCurrentPosition(scope.centerAt);
		}
	};
});

