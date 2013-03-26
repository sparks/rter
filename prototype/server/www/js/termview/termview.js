angular.module('termview', [
	'ng',      //filers
	'ui',      //ui-sortable and map
	'items',   //ItemCache to load items into termview, various itemDialog services
	'taxonomy' //Rankings
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
				this.termViews.push({term: term, heading: "All Items", active: true});
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

.controller('TermViewCtrl', function($scope, $filter, ItemCache, UpdateItemDialog, CloseupItemDialog, TermViewRemote, TaxonomyRankingCache) {
	/* -- items and rankings  -- */

	$scope.rankingCache = new TaxonomyRankingCache($scope.term.Term);

	if($scope.term.Term === "" || $scope.term.Term === undefined) {
		$scope.ranking = [];
	} else {
		$scope.ranking = $scope.rankingCache.ranking;
	}

	$scope.items = ItemCache.items;

	$scope.filteredItems = $filter('filterByTerm')($scope.items, $scope.term.Term);
	$scope.rankedItems = $filter('orderByRanking')($scope.filteredItems, $scope.ranking);

	$scope.$watch('items', function() {
		$scope.filteredItems = $filter('filterByTerm')($scope.items, $scope.term.Term);
	}, true);

	$scope.$watch('filteredItems', function() {
		$scope.rankedItems = $filter('orderByRanking')($scope.filteredItems, $scope.ranking);
	}, true);

	$scope.$watch('ranking', function() {
		$scope.rankedItems = $filter('orderByRanking')($scope.filteredItems, $scope.ranking);
	}, true);

	$scope.$watch('rankedItems', function(a, b) {
		$scope.updateMarkers();
	}, true);

	$scope.dragCallback = function(a) {
		var newRanking = [];
		angular.forEach($scope.rankedItems, function(v) {
			newRanking.push(v.ID);
		});

		if($scope.term.Term !== "" && $scope.term.Term !== undefined) {
			$scope.rankingCache.update(newRanking);
		}
	};

	$scope.closeupItemDialog = function(item){
		CloseupItemDialog.open(item);
	};

	$scope.updateItemDialog = function(item){
		UpdateItemDialog.open(item).then(function() {
			$scope.updateMarkers();
		});
	};

	$scope.close = function() {
		TermViewRemote.removeTermView($scope.term);
	};

	/* -- Map -- */

	$scope.markerBundles = [];

	$scope.mapCenter = new google.maps.LatLng(45.50745, -73.5793);

	$scope.mapOptions = {
		center: $scope.mapCenter,
		zoom: 10,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	$scope.mapResized = false;

	$scope.resizeMap = function() { //FIXME: Another map hack to render hidden maps
		if(!$scope.mapResized) {
			$scope.mapResized = true;
			google.maps.event.trigger($scope.map, "resize");
			$scope.map.setCenter($scope.mapCenter);
		}
	};

	$scope.updateMarkers = function() {
		angular.forEach($scope.markerBundles, function(v) {
			v.marker.setMap(null);
		});

		$scope.markerBundles = [];

		angular.forEach($scope.rankedItems, function(v) {
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

	$scope.centerAt = function(location) {
		var latlng = new google.maps.LatLng(location.coords.latitude, location.coords.longitude);
		$scope.map.setCenter(latlng);
		$scope.mapCenter = latlng;
	};
})

.directive('termview', function() {
	return {
		restrict: 'E',
		scope: {
			term: "="
		},
		templateUrl: '/template/termview/termview.html',
		controller: 'TermViewCtrl',
		link: function(scope, element, attrs) {
			navigator.geolocation.getCurrentPosition(scope.centerAt);
		}
	};
});

