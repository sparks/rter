function LivethumbnailDemoCtrl($scope, $http, $log) {
    $scope.videos = [];
    $http.get('assets/live-videos.json').success(function(data) {
            $scope.livevideos4thumbs = data.slice(0,6);
    });

	// test event reception
	$scope.$on("selected", function(e, video) {
		$log.info("VideoThumbnail: selected " + video.title);
	});

	$scope.$on("deselected", function(e, video) {
		$log.info("VideoThumbnail: deselected " + video.title);
	});

	$scope.$on("clicked", function(e, video) {
		$log.info("VideoThumbnail: clicked " + video.title);
	});

	$scope.$on("playing", function(e, video) {
		$log.info("VideoThumbnail: playing " + video.title);
	});

	$scope.$on("paused", function(e, video) {
		$log.info("VideoThumbnail: paused " + video.title);
	});

	$scope.$on("skimming", function(e, video) {
		$log.info("VideoThumbnail: skimming " + video.title);
	});

	$scope.$on("eos", function(e, video) {
		$log.info("VideoThumbnail: end-of-stream " + video.title);
	});

}

LivethumbnailDemoCtrl.$inject = ['$scope', '$http', '$log'];
