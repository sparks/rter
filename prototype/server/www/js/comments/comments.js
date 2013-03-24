angular.module('comments', [
	'ngResource',   //$resource for Item
	'alerts'       //Alerts for item actions
])

.factory('CommentRessource', function ($resource) {
	var CommentRessource = $resource(
		'/1.0/items/:ID/comments',
		{},
		{}
	);

	return CommentRessource;
})

.controller('CommentsDialogCtrl', function($scope, Alerter, CommentRessource) {
	$scope.comments = [];
	$scope.newComment = {};

	$scope.createComment = function() {

	};

})

.directive('commentsDialog', function() {
	return {
		restrict: 'E',
		scope: {},
		templateUrl: '/template/comments/comments-dialog.html',
		controller: 'CommentsDialogCtrl',
		link: function(scope, element, attrs) {

		}
	};
});