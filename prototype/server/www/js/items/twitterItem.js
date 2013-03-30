angular.module('twitterItem',  [
	'ng',   		//$timeout
	'ui',           //Map
	'ui.bootstrap'
])

.controller('FormTwitterItemCtrl', function($scope, $resource) {

	$scope.extra = {};
	$scope.extra.ResultType = "recent";
	
	if($scope.item.Author === undefined) {
		$scope.item.Author = "anonymous"; //TODO: Replace with login
	}
	$scope.mapCenter = new google.maps.LatLng(45.50745, -73.5793);

	$scope.mapOptions = {
		center: $scope.mapCenter,
		zoom: 10,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};

	$scope.centerAt = function(location) {
		var latlng = new google.maps.LatLng(location.coords.latitude, location.coords.longitude);
		$scope.map.setCenter(latlng);
		$scope.mapCenter = latlng;
	};

	$scope.setMarker = function($event) {
		if($scope.marker === undefined) {
			$scope.marker = new google.maps.Marker({
				map: $scope.map,
				position: $event.latLng
			});
		} else {
			$scope.marker.setPosition($event.latLng);
		}

		$scope.item.Lat = $event.latLng.lat();
		$scope.item.Lng = $event.latLng.lng();
	};


	
})
                                                                                                                                     
.directive('formTwitterItem', function($timeout) {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/twitter/form-twitter-item.html',
		controller: 'FormTwitterItemCtrl',
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
				
				scope.buildURL = function(){
						
					if(scope.extra.SearchTerm == undefined) {
						console.log("Error:  Search Query not set");
					}
						
					var searchURL = "http://search.twitter.com/search.json?page=1&rpp=40&callback=JSON_CALLBACK"	
										+ "&q=" + scope.extra.SearchTerm 
										+ "&result_type=" + scope.extra.ResultType

					if(!(scope.item.Lat == undefined)){
						searchURL = searchURL + "&geocode="+scope.item.Lat+","+scope.item.Lng+","+10+"mi";		
					}
										
					scope.item.ContentURI = encodeURI(searchURL);
					console.log("Built ContentURI " + scope.item.ContentURI);
				};


				$timeout( //FIXME: Another map hack to render hidden maps
					function() {
						google.maps.event.trigger(scope.map, "resize");
						scope.map.setCenter(scope.mapCenter);
					},
					0
				);

				console.log(scope.item, scope.extra);			
				
				scope.$watch('extra.SearchTerm', function(newVal, oldVal){
					if(!(newVal == undefined))	scope.buildURL();
				});
				scope.$watch('extra.ResultType', function(newVal, oldVal){
					if(!(newVal == undefined))	scope.buildURL();
				}, true);
				scope.$watch('item.Lat', function(newVal, oldVal){
					if(!(newVal == undefined))	scope.buildURL();
				}, true);

		}
	};
})

.controller('TileTwitterItemCtrl', function($scope) {

})

.directive('tileTwitterItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "="
		},
		templateUrl: '/template/items/twitter/tile-twitter-item.html',
		controller: 'TileTwitterItemCtrl',
		link: function(scope, element, attr) {
			
			
		}
	};
})
.controller('embedTweetCardCtrl', function($scope) {

})

.directive('embedTweetCard', function() {
	return {
		restrict: 'E',
		scope: {
			tweet: "="
		},
		templateUrl: '/template/items/twitter/twitter-card.html',
		controller: 'embedTweetCardCtrl',
		link: function(scope, element, attr) {
			
			
		}
	};
})

.controller('CloseupTwitterItemCtrl', function($scope, $http, ItemCache, CloseupItemDialog) {
	 console.log($scope.item.ContentURI);
	 $http({method: 'jsonp', url:$scope.item.ContentURI, cache: false}).
      success(function(data, status) {
          console.log(data, status);
        $scope.searchResult = data;
      }).
      error(function(data, status) {
         console.log(data, status);
        $scope.data = data || "Request failed";
        $scope.status = status;
    });
    
    /*
    $scope.showTweetCard = function(id, $event){
		// alert(id, $event.target);
		console.log(id, $event, $event.target);
		var urlVar = 'http://api.twitter.com/1/statuses/oembed.json?id='+id
					+'&align=center&omit_script=true&hide_thread=true&hide_media=true&callback=JSON_CALLBACK'
		$http({method: 'jsonp', url: urlVar, cache: false}).
	      success(function(data, status) {
	          console.log(data.html, status);
	          console.log($scope);
	        $scope.displayTweet =  data.html;
	        console.log($scope.displayTweet);

	      }).
	      error(function(data, status) {
	         console.log(data, status);
	        $scope.data = data || "Request failed";
	        $scope.status = status;
	    });
	};*/


	$scope.test = function(tweet, $event) {
			
			var newItem = {} ;
			newItem.Type = "SingleTweet";
			newItem.ContentURI = "http://twitter.com/{{tweet.from_user}}/status/{{tweet.id_str}}";

			console.log("it worked",$event );
			console.log("item - ",newItem );			
			ItemCache.create(
			{Type: "generic", ContentURI: tweet.id_str },
			function() {
			if($scope.dialog !== undefined) {
			$scope.dialog.close();
			}
		},
		function(e) {
			console.log(e);
		});
	};  
})

.directive('closeupTwitterItem', function() {
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