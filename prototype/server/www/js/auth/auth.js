angular.module('auth', [
	'ng',                    //$http
	'ui.bootstrap',          //dialog
	'http-auth-interceptor', //401 capture
	'ngResource',            //$resource
	'cache',                 //CacheBuilder
	'sockjs',                //sock for comment cache
	'alerts'                 //Alerter
])

.config(function(authServiceProvider) {
	authServiceProvider.addIgnoreUrlExpression(function (response) {
		return response.config.url === "/auth";
	});
})

.factory('UserResource', function ($resource) {
	var UserResource = $resource(
		'/1.0/users/:Username',
		{},
		{}
	);

	return UserResource;
})

.factory('UserDirectionResource', function ($resource) {
	var UserDirectionResource = $resource(
		'/1.0/users/:Username/direction',
		{ Username: '@Username' },
		{
			update: { method: 'PUT' }
		}
	);

	return UserDirectionResource;
})

.factory('UserDirectionCache', function($rootScope, SockJS, UserDirectionResource) {
	function UserDirectionCache(username) {
		var self = this;

		this.direction = {
			Username: username,
			Heading: 0
		};

		this.stream = new SockJS('/1.0/streaming/users/'+this.direction.Username+'/direction');

		function updateDirection(newDirection) {
			for(var key in newDirection) {
				self.direction[key] = newDirection[key];
			}
		}

		this.stream.onopen = function() {

		};

		this.stream.onmessage = function(e) {
			var bundle = e.data;

			if(bundle.Action == "update") {
				//Often if the user created the item, it will already be in place so treat as an update

				updateDirection(bundle.Val);
			}

			$rootScope.$digest();
		};

		this.stream.onclose = function() {

		};

		this.init = function() {
			UserDirectionResource.get(
				this.direction,
				function(direction) {
					updateDirection(direction);
				},
				function(e) {
					console.log(e);
				}
			);
		};

		this.close = function() {
			this.stream.close();
		}

		this.init();

		this.update = function(direction, sucess, failure) {
			UserDirectionResource.update(
				direction,
				function() {
					//Success do nothing!
					if(angular.isFunction(sucess)) sucess();
				},
				function(e) {
					if(e.status != 304) {
						Alerter.error("There was a problem updating the ranking. "+"Status:"+e.status+". Reply Body:"+e.data);
						console.log(e);
					}

					if(angular.isFunction(failure)) failure(e);
				}
			);
		};
	}

	return UserDirectionCache;
})

.controller('LoginPanelCtrl', function($scope, $http, authService, UserResource, Alerter) {
	$scope.login = function() {
		console.log("asdf");
		$http.post("/auth", {Username: $scope.username, Password: $scope.password})
		.success(function(data, status, headers) {
			$scope.cancel('success');
			authService.loginConfirmed();
		})
		.error(function(data, status, headers) {
			$scope.failedLogin = true;
			console.log("Login Problem", data, status);
			Alerter.error("Invalid login credentials.", 2000);
		});
	};

	$scope.signup = function() {
		UserResource.save(
			{Username: $scope.username, Password: $scope.password},
			function() {
				Alerter.success("User "+$scope.username+" created!", 2000);
			},
			function(e) {
				Alerter.error("Couldn't signup, Username already taken.", 2000);
				console.log(e);
			}
		);
	};

	$scope.cancel = function(result) {
		if($scope.dialog !== undefined) {
			$scope.dialog.close(result);
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
				dialogClass: 'modal login-modal',
				resolve: {item: function() { return item; }},
				templateUrl: '/template/auth/login-panel-dialog.html',
				controller: 'LoginDialogCtrl'
			});

			return d.open();
		}
	};
});