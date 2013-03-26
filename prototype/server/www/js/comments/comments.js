angular.module('comments', [
	'ngResource', //$resource for Item
	'alerts',     //Alerts for item actions
	'moment'     //fromNow filter
])

.factory('CommentResource', function ($resource) {
	var CommentResource = $resource(
		'/1.0/items/:ID/comments',
		{ ID: '@ID' },
		{}
	);

	return CommentResource;
})

.controller('CommentsDialogCtrl', function($scope, Alerter, CommentResource) {
	$scope.comments = CommentResource.query({ID: $scope.id});

	$scope.newComment = {
		ID: $scope.id,
		Body: "",
		Author: "anonymous"
	};

	$scope.createComment = function() {
		if($scope.newComment.Body === undefined || $scope.newComment.Body === "") return;
		CommentResource.save(
			$scope.newComment,
			function(c) {
				$scope.comments.push(c);
				$scope.newComment = {
					ID: $scope.id,
					Body: "",
					Author: "anonymous"
				};
			},
			function(e) {
				console.log(e);
			}
		);
	};
})

.directive('commentsDialog', function(CommentResource) {
	return {
		restrict: 'E',
		scope: {
			id: "="
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
