angular.module('youtubeItem', [
	'ng',       //$timeout
	'edit-map', //Maps
	'disp-map'  //Maps
])

.controller('FormYoutubeItemCtrl', function($scope) {
	$scope.item.StartTime = new Date();
	$scope.item.StopTime = $scope.item.StartTime;

	$scope.item.HasHeading = false;
	$scope.item.HasGeo = false;
	$scope.item.Live = false;
})

.directive('formYoutubeItem', function($timeout) {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/youtube/form-youtube-item.html',
		controller: 'FormYoutubeItemCtrl',
		link: function(scope, element, attr) {

		}
	};
})

.controller('TileYoutubeItemCtrl', function($scope) {
	$scope.youtubeID = $scope.item.ContentURI.match(/\/watch\?v=([0-9a-zA-Z].*)/)[1];
})

.directive('tileYoutubeItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/youtube/tile-youtube-item.html',
		controller: 'TileYoutubeItemCtrl',
		link: function(scope, element, attr) {

		}
	};
})

.controller('CloseupYoutubeItemCtrl', function($scope) {
	$scope.youtubeID = $scope.item.ContentURI.match(/\/watch\?v=([0-9a-zA-Z].*)/)[1];
})

.directive('closeupYoutubeItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/youtube/closeup-youtube-item.html',
		controller: 'CloseupYoutubeItemCtrl',
		link: function(scope, element, attr) {

		}
	};
});