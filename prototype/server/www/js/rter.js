angular.module('rter', ['ui.bootstrap', 'items', 'termview', 'alerts', 'taxonomy'])

.controller('TabsCtrl', function($scope, Taxonomy) {
	$scope.termViews = [];
	$scope.terms = Taxonomy.query();

	$scope.addTermView = function(term) {
		for(var i = 0;i < $scope.termViews.length;i++) {
			if($scope.termViews[i].term.Term == term.Term) {
				$scope.termViews[i].active = true;
				return;
			}
		}

		$scope.termViews.push({term: term, active: true});
	};
})

.directive('eatClick', function() {
    return function(scope, element, attrs) {
        $(element).click(function(event) {
            event.preventDefault();
        });
    };
})

.directive('mapResize', function() {
	return function(scope, element, attrs) {
		$scope.$watch("mapVisible", function (v) {
			if (v) {
				google.maps.event.trigger($scope.map, "resize");
			}
		});
	};
});
