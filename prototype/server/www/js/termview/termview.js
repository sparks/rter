angular.module('termview', ['ngResource', 'items'])

.controller('TermViewCtrl', function($scope, Item) {
})

.directive('termview', function(Item) {
	return {
		restrict: 'E',
		scope: {
			term: "@"
		},
		templateUrl: '/template/termview/termview.html',
		controller: 'TermViewCtrl',
		link: function(scope, element, attrs) {
			if(attrs.term === undefined) {
				scope.items = Item.query();
			} else {
				scope.items = Item.query({term: attrs.term});
			}
		}
	};
});

