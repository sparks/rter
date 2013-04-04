angular.module('comments', [
	'ngResource', //$resource for Item
	'alerts',     //Alerts for item actions
	'cache',      //CacheBuilder
	'sockjs',     //sock for comment cache
	'moment'      //fromNow filter
])

.factory('CommentResource', function ($resource) {
	var CommentResource = $resource(
		'/1.0/items/:ItemID/comments',
		{ ItemID: '@ItemID' },
		{}
	);

	return CommentResource;
})

.factory('CommentCacheBuilder', function (CacheBuilder, $resource, SockJS) {
	return function(itemID) {
		return new CacheBuilder(
			"Comment",
			$resource(
				'/1.0/items/'+itemID+'/comments/:ID',
				{ ID: '@ID' },
				{}
			),
			new SockJS('/1.0/streaming/items/'+itemID+'/comments'),
			function(a, b) {
				if(a.ID === undefined || b.ID === undefined) return false;
				if(a.ID == b.ID) return true;
				return false;
			}
		);
	};
})
.controller('CommentsDialogCtrl', function($scope, Alerter, CommentCacheBuilder) {
	$scope.commentCache = CommentCacheBuilder($scope.itemId);
	$scope.comments = $scope.commentCache.contents;

	$scope.newComment = {
		Body: ""
	};

	$scope.inProgress = false;

	$scope.createComment = function() {
		if($scope.newComment.Body === undefined || $scope.newComment.Body === "") return;
		$scope.inProgress = true;
		$scope.commentCache.create(
			$scope.newComment,
			function(c) {
				$scope.newComment = {
					Body: ""
				};
				$scope.inProgress = false;
			},
			function(e) {
				console.log(e);
				$scope.inProgress = false;
			}
		);
	};

	$scope.$on("$destroy", function() {
		$scope.commentCache.close();
	});
})

.directive('commentsDialog', function() {
	return {
		restrict: 'E',
		scope: {
			itemId: "="
		},
		templateUrl: '/template/comments/comments-dialog.html',
		controller: 'CommentsDialogCtrl',
		link: function(scope, element, attrs) {
			scope.$watch(
				'comments',
				function() {
					element.children().children().scrollTop(0);
				},
				true
			);
		}
	};
});
