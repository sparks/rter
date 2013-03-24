angular.module('comments', [
	'ngResource',   //$resource for Item
	'alerts'       //Alerts for item actions
])

.factory('CommentRessource', function ($resource) {
	var CommentRessource = $resource(
		'/1.0/items/:ID/comments',
		{},
		{
			save: { method: 'POST', params:{ ID: '@ID' } },
			query: { method: 'GET', params:{ ID: '@ID' }, isArray: true }
		}
	);

	return CommentRessource;
})

.controller('CommentsDialogCtrl', function($scope, Alerter, CommentRessource) {
	$scope.comments = CommentRessource.query({ID: $scope.id});

	$scope.newComment = {
		ID: $scope.id,
		Body: "",
		Author: "anonymous"
	};

	$scope.createComment = function() {
		CommentRessource.save(
			$scope.newComment,
			function() {
				$scope.comments.push($scope.newComment);
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

.directive('commentsDialog', function(CommentRessource) {
	return {
		restrict: 'E',
		scope: {
			id: "="
		},
		templateUrl: '/template/comments/comments-dialog.html',
		controller: 'CommentsDialogCtrl',
		link: function(scope, element, attrs) {
		}
	};
});
