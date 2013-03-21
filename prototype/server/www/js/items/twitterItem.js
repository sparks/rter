angular.module('twitterItem',  [
	'ui',           //Map
	'ui.bootstrap', //select2
	'taxonomy'      //Tag list
])

.controller('FormTwitterItemCtrl', function($scope) {

})

.directive('formTwitterItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/twitter/form-twitter-item.html',
		controller: 'FormTwitterItemCtrl',
		link: function(scope, element, attr) {

		}
	};
})

.controller('TileTwitterItemCtrl', function($scope) {

})

.directive('tileTwitterItem', function(Taxonomy) {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/twitter/tile-twitter-item.html',
		controller: 'TileTwitterItemCtrl',
		link: function(scope, element, attr) {
			if(scope.item.Lat !== undefined && scope.item.Lng !== undefined) {
				var latLng = new google.maps.LatLng(scope.item.Lat, scope.item.Lng);
				scope.marker = new google.maps.Marker({
					map: scope.map,
					position: latLng
				});
				scope.mapCenter = latLng;
			} else {
				navigator.geolocation.getCurrentPosition(scope.centerAt);
			}

			
		}
	};
})

.controller('CloseupTwitterItemCtrl', function($scope) {

})

.directive('closeupTwitterItem', function(Taxonomy) {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/twitter/closeup-twitter-item.html',
		controller: 'CloseupTwitterItemCtrl',
		link: function(scope, element, attr) {

		}
	};
});