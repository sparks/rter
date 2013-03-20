angular.module('rter', ['ui.bootstrap', 'items', 'termview', 'alerts', 'taxonomy'])

.controller('TabsCtrl', function($scope, Taxonomy, TermViewRemote) {
	$scope.termViews = TermViewRemote.termViews;

	TermViewRemote.addTermView({Term: ""});
	$scope.terms = Taxonomy.query();

	$scope.addTermView = TermViewRemote.addTermView;
})

.directive('eatClick', function() {
    return function(scope, element, attrs) {
        $(element).click(function(event) {
            event.preventDefault();
        });
    };
});