angular.module('singleItem', [
	'ng',   		//$timeout
	'ui',           //Map
	'ui.bootstrap'
])


.controller('TileSingleTweetItemCtrl', function($scope) {

})

.directive('tileSingleTweetItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/singletweet/tile-singletweet-item.html',
		controller: 'TileSingleTweeItemCtrl',
		link: function(scope, element, attr) {

		}
	};
})

.controller('CloseupSingleTweeItemCtrl', function($scope, $http) {
			console.log($scope.item.ContentURI);
			$http({method: 'jsonp', url: $scope.item.ContentURI, cache: false}).
		      success(function(data, status) {
		          console.log(status);
		          console.log($scope);
		        $scope.displayTweet =  data.html;
		        var TweetCardHtml = angular.element(data.html);
		        $('#tweetcard').append(TweetCardHtml);
		        console.log($scope.displayTweet);

		      }).
		      error(function(data, status) {
		         console.log(data, status);
		        $scope.data = data || "Request failed";
		        $scope.status = status;
		    });
})

.directive('closeupSingletweetItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/singletweet/closeup-singletweet-item.html',
		controller: 'CloseupSingleTweeItemCtrl',
		link: function(scope, element, attr) {


		}
	};
});