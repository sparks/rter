angular.module('rter', ['ui.bootstrap', 'items', 'termview', 'alerts'])

.controller('TabsCtrl', function($scope) {
	$scope.termviews = [
		// {term: 'a'}
	];
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
