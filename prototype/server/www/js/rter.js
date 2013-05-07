angular.module('rter', [
	'ui.bootstrap',          //Tabs
	'alerts',                //Main alert box
	'taxonomy',              //Taxonomy for tag-cloud
	'termview',              //term-view directives and TermViewRemote
	'http-auth-interceptor', //401 catcher
	'ngCookies',             //Cookie for login/logout
	'auth',                  //Login system
	'tsunami'                //live video player
])

.filter('if', function() {
	return function(input, value) {
		if (typeof(input) === 'string') {
			input = [input, ''];
		}
		return value? input[0] : input[1];
	};
})

.controller('RterCtrl', function($scope, LoginDialog, $cookies, $cookieStore) {
	if($cookies['rter-credentials'] !== undefined) {
		$scope.loggedIn = true;
	} else {
		$scope.loggedIn = false;
	}

	$scope.loginDialogOpen = false;

	$scope.login = function() {
		if(!$scope.loginDialogOpen) {
			$scope.loginDialogOpen = true;
			LoginDialog.open().then(function(result) {
				if(result == 'success') {
					$scope.loggedIn = true;
				}
				$scope.loginDialogOpen = false;
			});
		}
	};

	$scope.logout = function() {
		$cookieStore.remove('rter-credentials');
		$scope.loggedIn = false;
	};

	$scope.$on('event:auth-loginRequired', function() {
		$scope.login();
	});
})

.controller('TabsCtrl', function($scope, TermViewRemote) {
	$scope.termViews = TermViewRemote.termViews;
	TermViewRemote.addTermView({Term: ""});
})

.controller('TagCloudCtrl', function($scope, TermViewRemote, TaxonomyCache) {
	$scope.terms = TaxonomyCache.contents //TODO: Make me dynamic

	$scope.$watch('terms', function() {
		$scope.countMax = 0;

		angular.forEach($scope.terms, function(val) {
			if($scope.countMax < val.Count) $scope.countMax = val.Count;
		});
	}, true);

	$scope.addTermView = function(term) {
		TermViewRemote.addTermView(term);
	};

	$scope.termFontSize = function(term) {
		return map(term.Count, 1, $scope.countMax, 12, 60);
	};
})

.directive('tagCloud', function(ItemCache) {
	return {
		restrict: 'E',
		scope: {},
		templateUrl: '/template/tag-cloud.html',
		controller: 'TagCloudCtrl',
		link: function(scope, element, attrs) {

		}
	};
});
