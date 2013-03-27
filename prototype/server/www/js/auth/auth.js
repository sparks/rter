angular.module('auth', [
	'ui.bootstrap',         //dialog
	'http-auth-interceptor' //$resource for taxonomoy
])

.factory('UserResource', function ($resource) {
	var UserResource = $resource(
		'/1.0/users/:Username',
		{},
		{}
	);

	return UserResource;
})

.controller('LoginPanelCtrl', function($scope, authService, UserResource) {
	$scope.login = function() {
		console.log($scope.username);
		console.log($scope.password);

		$scope.cancel();

		authService.loginConfirmed();
	};

	$scope.signup = function() {
		UserResource.save(
			{Username: $scope.username, Password: $scope.password},
			function() {
				console.log("done");
			},
			function(e) {
				console.log(e);
			}
		);
	};

	$scope.cancel = function() {
		if($scope.dialog !== undefined) {
			$scope.dialog.close();
		}
	};
})

.directive('loginPanel', function(authService) {
	return {
		restrict: 'E',
		scope: {
			dialog: "="
		},
		templateUrl: '/template/auth/login-panel.html',
		controller: 'LoginPanelCtrl',
		link: function(scope, element, attrs) {

		}
	};
})

.controller('LoginDialogCtrl', function($scope, dialog) {
	$scope.dialog = dialog;
})

.factory('LoginDialog', function ($dialog) {
	return {
		open: function(item) {
			var d = $dialog.dialog({
				modalFade: false,
				backdrop: true,
				keyboard: true,
				backdropClick: false,
				resolve: {item: function() { return item; }},
				templateUrl: '/template/auth/login-panel-dialog.html',
				controller: 'LoginDialogCtrl'
			});

			return d.open();
		}
	};
});