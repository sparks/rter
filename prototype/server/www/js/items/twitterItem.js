angular.module('twitterItem',  [
	'ngResource',   //Twitter Rest API
	'ui',           //Map
	'ui.bootstrap'
])

.controller('FormTwitterItemCtrl', function($scope, $resource) {

	if($scope.item.Author === undefined) {
		$scope.item.Author = "anonymous"; //TODO: Replace with login
	}

	//This is kinda terrible
	if($scope.item.Terms !== undefined) {
		var concat = "";
		for(var i = 0;i < $scope.item.Terms.length;i++) {
			concat += $scope.item.Terms[i].Term+",";
		}
		$scope.item.Terms = concat.substring(0, concat.length-1);
	}

	

	$scope.item.ContentURI = 'http://search.twitter.com/:action'+'search.json';
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
	};
})

.controller('CloseupTwitterItemCtrl', function($scope) {
	alert('We are in the closeupcontroller');
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