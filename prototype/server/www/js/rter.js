angular.module('rter', [
	'ui.bootstrap', //Tabs
	'alerts',       //Main alert box
	'taxonomy',     //Taxonomy for tag-cloud
	'termview'      //term-view directives and TermViewRemote
])

.controller('TabsCtrl', function($scope, TermViewRemote) {
	$scope.termViews = TermViewRemote.termViews;
	TermViewRemote.addTermView({Term: ""});
})

.directive('eatClick', function() {
    return function(scope, element, attrs) {
        $(element).click(function(event) {
            event.preventDefault();
        });
    };
})

.controller('TagCloudCtrl', function($scope, TermViewRemote, Taxonomy) {
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

.directive('tagCloud', function(ItemCache) {
	return {
		restrict: 'E',
		scope: {},
		templateUrl: '/template/tag-cloud.html',
		controller: 'TagCloudCtrl',
		link: function(scope, element, attrs) {

		}
	};
});