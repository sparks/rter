angular.module('twitterItem', [])

.controller('FormTwitterItemCtrl', function($scope) {

})

.directive('formTwitterItem', function(Item) {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/form-twitter-item.html',
		controller: 'FormTwitterItemCtrl',
		link: function(scope, element, attr) {

		}
	};
});