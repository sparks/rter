angular.module('streamingVideoV1Item', [
	'tsunamijs.livethumbnail', //Live Thumbnails
	'edit-map',                //maps
	'disp-map'                 //maps
])

.controller('TileStreamingVideoV1ItemCtrl', function($scope) {
	$scope.video = {};

	$scope.$watch('item', function() {
		if(!$scope.item) return;

		$scope.video.title = "Live";
		$scope.video.thumbnailUrl = $scope.item.ThumbnailURI;
		$scope.video.StartTime = $scope.item.StartTime;
		$scope.video.EndTime = $scope.item.StopTime;
	}, true);
})

.directive('tileStreamingVideoV1Item', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/streamingVideoV1/tile-streamingVideoV1-item.html',
		controller: 'TileStreamingVideoV1ItemCtrl',
		link: function(scope, element, attr) {

		}
	};
})

.controller('CloseupStreamingVideoV1ItemCtrl', function($scope, ItemCache, CloseupItemDialog) {

})

.directive('closeupStreamingVideoV1Item', function() {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			dialog: "="
		},
		templateUrl: '/template/items/streamingVideoV1/closeup-streamingVideoV1-item.html',
		controller: 'CloseupStreamingVideoV1ItemCtrl',
		link: function(scope, element, attr) {

		}
	};
})

.directive('ngPoster', function() {
	return {
		priority: 99, // it needs to run after the attributes are interpolated
		link: function(scope, element, attr) {
			attr.$observe('ngPoster', function(value) {
				if (!value)
					return;

				attr.$set('poster', value);
			});
		}
	};
})

.directive('autoplayIf', function() {
	return {
		priority: 99, // it needs to run after the attributes are interpolated
		link: function(scope, element, attr) {
			attr.$observe('autoplayIf', function(value) {
				if (!value)
					return;

				attr.$set('autoplay', '');
			});
		}
	};
});
