angular.module('rter', ['ui.bootstrap', 'items', 'termview', 'alerts', 'taxonomy'])

.controller('TabsCtrl', function($scope, Taxonomy, TermViewRemote) {
	$scope.termViews = TermViewRemote.termViews;

	TermViewRemote.addTermView({Term: ""});
	$scope.terms = Taxonomy.query(function() {
		$scope.countMax = 0;

		angular.forEach($scope.terms, function(val) {
			if($scope.countMax < val.Count) $scope.countMax = val.Count;
		});
	});

	$scope.addTermView = TermViewRemote.addTermView;

	$scope.termFontSize = function(term) {
		return term.Count/$scope.countMax*30;
	};
})

.directive('eatClick', function() {
    return function(scope, element, attrs) {
        $(element).click(function(event) {
            event.preventDefault();
        });
    };
});