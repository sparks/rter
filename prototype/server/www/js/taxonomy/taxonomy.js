angular.module('taxonomy', ['ngResource', 'termview'])

.factory('Taxonomy', function ($resource) {
	var Taxonomy = $resource(
		'/1.0/taxonomy/:Term',
		{},
		{
			update: { method: 'PUT', params:{ Term: '@Term' } }
		}
	);

	return Taxonomy;
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
		templateUrl: '/template/taxonomy/tag-cloud.html',
		controller: 'TagCloudCtrl',
		link: function(scope, element, attrs) {

		}
	};
});