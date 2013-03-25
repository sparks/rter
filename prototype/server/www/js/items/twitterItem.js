angular.module('twitterItem',  [
	'ngResource',   //Twitter Rest API
	'ui',           //Map
	'ui.bootstrap'
])

.controller('FormTwitterItemCtrl', function($scope, $resource) {

	console.log("Inside Twitter Form Ctrl");
	if($scope.item.Author === undefined) {
		$scope.item.Author = "anonymous"; //TODO: Replace with login
	}
	

	
})

.directive('formTwitterItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/twitter/form-twitter-item.html',
		controller: 'FormTwitterItemCtrl',
		link: function(scope, element, attr) {
			console.log(scope.item);			
			scope.$watch('item.searchTerm', function(newVal, oldVal){
				console.log("Inside the watch");
				scope.item.ContentURI = 'http://search.twitter.com/search.json?q='+newVal;
			})

			
		}
	};
})

.controller('TileTwitterItemCtrl', function($scope) {

})

.directive('tileTwitterItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/twitter/tile-twitter-item.html',
		controller: 'TileTwitterItemCtrl',
		link: function(scope, element, attr) {
			
			
		}
	}
;})

.controller('CloseupTwitterItemCtrl', function($scope) {
	
	$scope.twitterConfig = $resource('http://search.twitter.com/:action',
		{action: 'search.json', q:'montreal', callback: 'JSON_CALLBACK'},	
		{get:{method : 'JSONP'}}
	);

	$scope.twitterConfig.get();
})

.directive('closeupTwitterItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/twitter/closeup-twitter-item.html',
		controller: 'CloseupTwitterItemCtrl',
		link: function(scope, element, attr) {

		}
	};
});