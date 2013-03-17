angular.module('termview', ['ngResource', 'items', 'ui.bootstrap.dialog'])

.controller('TermViewCtrl', function($scope, $dialog, Item) {
	$scope.updateItemDialog = function(item){
		var d = $dialog.dialog(
			{
				modalFade: false,
				backdrop: true,
				keyboard: true,
				backdropClick: true,
				resolve: {item: function() { return item; }},
				templateUrl: '/template/items/submit-item.html',
				controller: 'UpdateItemCtrl'
			}
		);
		d.open();
	};
})

.directive('termview', function(Item) {
	return {
		restrict: 'E',
		scope: {
			term: "@"
		},
		templateUrl: '/template/termview/termview.html',
		controller: 'TermViewCtrl',
		link: function(scope, element, attrs) {
			if(attrs.term === undefined) {
				scope.items = Item.query();
			} else {
				scope.items = Item.query({term: attrs.term});
			}
		}
	};
});

