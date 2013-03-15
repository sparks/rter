angular.module('twitterItem', [])

.directive('submitTwitterItem', function(Item) {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/submit-twitter-item.html',
		link: function(scope, element, attr) {

		}
	};
});