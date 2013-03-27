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
	

	// 871025159  DELL kate leave here a message
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

.controller('CloseupTwitterItemCtrl', function($scope, $http) {

	 $http({method: 'jsonp', url: 'https://api.twitter.com/1/statuses/oembed.json?id=316661563513782272&align=center&callback=JSON_CALLBACK', cache: false}).
      success(function(data, status) {
          console.log(data.html, status);
          console.log($scope);
        $scope.displayTweet =  data.html;
        console.log($scope.displayTweet);
        $scope.status = status;
        $scope.data = data;
      }).
      error(function(data, status) {
         console.log(data, status);
        $scope.data = data || "Request failed";
        $scope.status = status;
    });
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
			scope.addHTML = function(newhtml) {
				element.append($(newhtl))
			}

		}
	};
});