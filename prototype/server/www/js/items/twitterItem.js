angular.module('twitterItem', [])

.controller('SubmitTwitterItemCtrl', function($scope) {

})

.directive('submitTwitterItem', function(Item) {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/submit-twitter-item.html',
		controller: 'SubmitTwitterItemCtrl',
		link: function(scope, element, attr) {

		}
	};
});