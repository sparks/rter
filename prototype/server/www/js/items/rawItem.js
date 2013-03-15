angular.module('rawItem', ['ui.directives'])

.directive('submitRawItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/submit-raw-item.html',
		link: function(scope, element, attr) {

		}
	};
});