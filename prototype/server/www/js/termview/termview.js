angular.module('termview', [
	'ng',       //filers
	'ui',       //ui-sortable and map
	'items',    //ItemCache to load items into termview, various itemDialog services
	'taxonomy', //Rankings
	'alerts'    //Alerter
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

.controller('TermViewCtrl', function($scope, $filter, $timeout, Alerter, ItemCache, UpdateItemDialog, CloseupItemDialog, TermViewRemote, TaxonomyRankingCache) {

	$scope.viewmode = "grid-view";
	$scope.filterMode = "blur";
	$scope.prevFilterMode = "blur";
	$scope.mapFilterEnable = false;

	$scope.$watch('viewmode', function(newVal, oldVal) {
		$scope.mapCenter = $scope.map.getCenter();

		if($scope.viewmode == 'map-view') {
			$scope.mapFilterEnable = false;
		}

		if(oldVal != "map-view" && newVal == "map-view") {
			$scope.prevFilterMode = $scope.filterMode;
			$scope.filterMode = "remove";
		}

		if(oldVal == "map-view" && newVal != "map-view") {
			$scope.filterMode = $scope.prevFilterMode;
		}

		$timeout(function() {
			$scope.resizeMap();
		}, 0);
	});

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

	$scope.finalMapItems = $scope.rankedItems;
	$scope.finalFilteredItems = $scope.rankedItems;

	$scope.textSearchedItems = $filter('filter')($scope.rankedItems, $scope.filterQuery);
	$scope.mapFilteredItems = $filter('filterbyBounds')($scope.textSearchedItems, $scope.mapBounds);

	$scope.$watch('items', function() {
		$scope.filteredItems = $filter('filterByTerm')($scope.items, $scope.term.Term);
	}, true);

	$scope.$watch('[ranking, filteredItems]', function() {
		$scope.rankedItems = $filter('orderByRanking')($scope.filteredItems, $scope.ranking);
	}, true);

	$scope.$watch('[rankedItems, filterMode]', function() {
		if($scope.filterMode == 'blur') {
			$scope.updateMarkers();
			$scope.finalFilteredItems = $scope.rankedItems;
			$scope.finalMapItems = $scope.rankedItems;
		}
	}, true);

	$scope.$watch('[rankedItems, textQuery, filterMode]', function() {
		if($scope.filterMode == 'remove') {
			$scope.textSearchedItems = $filter('filter')($scope.rankedItems, $scope.textQuery);
		}
	}, true);

	$scope.$watch('[textSearchedItems, filterMode]', function() {
		if($scope.filterMode == 'remove') {
			$scope.finalMapItems = $scope.textSearchedItems;
			$scope.updateMarkers();
		}
	}, true);

	$scope.$watch('[textSearchedItems, mapBounds, mapFilterEnable, filterMode]', function() {
		if($scope.filterMode == 'remove') {
			if($scope.mapFilterEnable) {
				$scope.mapFilteredItems = $filter('filterbyBounds')($scope.textSearchedItems, $scope.mapBounds);
			} else {
				$scope.mapFilteredItems = $scope.textSearchedItems;
			}
		}
	}, true);

	$scope.$watch('[mapFilteredItems, filterMode]', function() {
		if($scope.filterMode == 'remove') {
			$scope.finalFilteredItems = $scope.mapFilteredItems;
		}
	}, true);

	$scope.isFiltered = function(item) {
		var filtered = [item];

		filtered = $filter('filter')(filtered, $scope.textQuery);

		if($scope.mapFilterEnable) {
			filtered = $filter('filterbyBounds')(filtered, $scope.mapBounds);
		}

		if(filtered.length === 0) return true;
		else return false;
	};

	$scope.dragCallback = function(a, b) {
		if($scope.filterMode == 'remove' && ($scope.mapFilterEnable || ($scope.textQuery !== undefined && $scope.textQuery !== ''))) { //TODO: This should have a blur options instead maybe?
			Alerter.warn("You cannot reorder items while your filters are enabled", 2000);
			return;
		}

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

	$scope.boundsChanged = function() {
		$scope.mapBounds = $scope.map.getBounds();
	};

	$scope.markerBundles = [];

	$scope.mapCenter = new google.maps.LatLng(45.50745, -73.5793);

	$scope.mapOptions = {
		center: $scope.mapCenter,
		zoom: 10,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	$scope.resizeMap = function() {
		google.maps.event.trigger($scope.map, "resize");
		$scope.map.setCenter($scope.mapCenter);
		$scope.mapBounds = $scope.map.getBounds();
	};

	$scope.updateMarkers = function() {
		angular.forEach($scope.markerBundles, function(v) {
			v.marker.setMap(null);
		});

		$scope.markerBundles = [];

		angular.forEach($scope.finalMapItems, function(v) {
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

