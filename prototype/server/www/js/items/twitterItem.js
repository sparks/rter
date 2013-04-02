angular.module('twitterItem',  [
	'ng',   		//$timeout
	'ui',           //Map
	'ui.bootstrap'
])

.controller('FormTwitterItemCtrl', function($scope, $resource) {

	$scope.extra = {};
	$scope.extra.ResultType = "recent";
	
	$scope.mapCenter = new google.maps.LatLng(45.50745, -73.5793);

	$scope.mapOptions = {
		center: $scope.mapCenter,
		zoom: 15,
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
			//making the circle
			$scope.extra.radius = 10*1000; //10km in meters
			$scope.circle = new google.maps.Circle({
				map: $scope.map,
				center: $scope.marker.getPosition(),
				radius: $scope.extra.radius,
				editable: true,
				draggable: false,
				fillColor: "#FF0000",
				fillOpacity: 0.3,
				strokeColor: "#FF0000",
			    strokeOpacity: 0.8,
			    strokeWeight: 2

			});
			google.maps.event.addListener($scope.circle, 'radius_changed', function() {
				$scope.extra.radius = $scope.circle.getRadius();
				console.log("Radius changed and calling buildURl");
				$scope.buildURL();
			});
			google.maps.event.addListener($scope.circle, 'center_changed', function() {
			  	$scope.marker.setPosition($scope.circle.getCenter());
			  	console.log("Center changed and calling buildURl");
			  	$scope.buildURL();
			});
		} else {
			$scope.marker.setPosition($event.latLng);
			$scope.circle.setCenter($event.latLng);
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
						searchURL = searchURL + "&geocode="+scope.item.Lat+","+scope.item.Lng+","+(scope.extra.radius/1000)+"km";		
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
			newItem.ContentURI = "http://twitter.com/"+tweet.from_user+"/status/"+tweet.id_str;
			console.log("it worked",$event );
			if(tweet.geo != null)
			{
				newItem.Lat = tweet.geo.coordinates[0];
				newItem.Lng = tweet.geo.coordinates[1];					
			}
			console.log(newItem);
			
			ItemCache.create(
			newItem,
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