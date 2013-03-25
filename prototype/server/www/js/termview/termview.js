angular.module('termview', [
	'ng',      //filers
	'ui',      //ui-sortable and map
	'items',   //ItemCache to load items into termview, various itemDialog services
	'taxonomy' //
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

.controller('TermViewCtrl', function($scope, $filter, UpdateItemDialog, CloseupItemDialog, TermViewRemote, TaxonomyRanking) {
	$scope.ranking = TaxonomyRanking.get(
		{Term: $scope.term.Term},
		function() {
			console.log($scope.ranking);
			console.log($scope.ranking.Ranking);
			console.log(JSON.parse($scope.ranking.Ranking));
		}
	);

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

	$scope.markerBundles = [];

	$scope.updateMarkers = function() {
		var filteredItems = $filter('filterByTerm')($scope.items, $scope.term.Term);

		angular.forEach($scope.markerBundles, function(v) {
			v.marker.setMap(null);
		});

		$scope.markerBundles = [];

		angular.forEach(filteredItems, function(v) {
			if(v.Lat === undefined || v.Lng === undefined || (v.Lat === 0 && v.Lng === 0)) return;

			var m = new google.maps.Marker({
				map: $scope.map,
				position: new google.maps.LatLng(v.Lat, v.Lng)
			});

			if(v.ThumbnailURI !== undefined && v.ThumbnailURI !== "") {
				m.setIcon(new google.maps.MarkerImage(v.ThumbnailURI, null, null, null, new google.maps.Size(40, 40)));
			}

			$scope.markerBundles.push({marker: m, item: v});
		});
	};

	$scope.openMarkerInfo = function(m) {
		console.log(m);
	};

	$scope.centerAt = function(location) {
		var latlng = new google.maps.LatLng(location.coords.latitude, location.coords.longitude);
		$scope.map.setCenter(latlng);
		$scope.mapCenter = latlng;
	};

	$scope.dragCallback = function() {
		console.log("whamo");
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

