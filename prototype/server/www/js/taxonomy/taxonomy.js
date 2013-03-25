angular.module('taxonomy', [
	'ngResource' //$resource for taxonomoy
])

.factory('TaxonomyRankingResource', function ($resource) {
	var TaxonomyRankingResource = $resource(
		'/1.0/taxonomy/:Term/ranking',
		{ Term: '@Term' },
		{
			update: { method: 'PUT', params:{ Term: '@Term' } }
		}
	);

	return TaxonomyRankingResource;
})

.factory('TaxonomyRankingCache', function($rootScope, TaxonomyRankingResource) {

})

.factory('TaxonomyResource', function ($resource) {
	var TaxonomyResource = $resource(
		'/1.0/taxonomy/:Term',
		{ Term: '@Term' },
		{
			update: { method: 'PUT', params:{ Term: '@Term' } }
		}
	);

	return TaxonomyResource;
})

.controller('TagSelectorCtrl', function($scope, TaxonomyResource) {
	if($scope.terms !== undefined) {
		var concat = "";
		for(var i = 0;i < $scope.terms.length;i++) {
			concat += $scope.terms[i].Term+",";
		}
		$scope.terms = concat.substring(0, concat.length-1);
	}

	$scope.tagConfig = {
		data: TaxonomyResource.query(),
		multiple: true,
		id: function(item) {
			return item.Term;
		},
		formatResult: function(item) {
			return item.Term;
		},
		formatSelection: function(item) {
			return item.Term;
		},
		createSearchChoice: function(term) {
			return {Term: term};
		},
		matcher: function(term, text, option) {
			return option.Term.toUpperCase().indexOf(term.toUpperCase())>=0;
		},
		initSelection: function (element, callback) {
			var data = [];
			$(element.val().split(",")).each(function (v, a) {
				data.push({Term: a});
			});
			callback(data);
		}
	};
})

.directive('tagSelector', function() {
	return {
		restrict: 'E',
		scope: {
			terms: "="
		},
		templateUrl: '/template/taxonomy/tag-selector.html',
		controller: 'TagSelectorCtrl',
		link: function(scope, element, attrs) {

		}
	};
});