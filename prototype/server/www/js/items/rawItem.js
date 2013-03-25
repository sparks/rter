angular.module('rawItem', [
	'ui',           //Map
	'ui.bootstrap', //select2
])

.controller('FormRawItemCtrl', function($scope) {

})

.directive('formRawItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/raw/form-raw-item.html',
		controller: 'FormRawItemCtrl',
		link: function(scope, element, attr) {

		}
	};
})

.controller('CloseupRawItemCtrl', function($scope) {

})

.directive('closeupRawItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/raw/closeup-raw-item.html',
		controller: 'CloseupRawItemCtrl',
		link: function(scope, element, attr) {

		}
	};
});