angular.module('taxonomy', [
	'ngResource', //$resource for taxonomoy
	'sockjs' //sock for TaxonomyRankingCache
])

.factory('TaxonomyRankingResource', function ($resource) {
	var TaxonomyRankingResource = $resource(
		'/1.0/taxonomy/:Term/ranking',
		{ Term: '@Term' },
		{
			update: { method: 'PUT' }
		}
	);

	return TaxonomyRankingResource;
})

.factory('TaxonomyRankingCache', function($rootScope, SockJS, TaxonomyRankingResource) {
	function TaxonomyRankingCache(term) {
		if(term === "" || term === undefined) return;

		var self = this;

		this.term = term;
		this.ranking = [];
		this.stream = new SockJS('/1.0/streaming/taxonomy/'+term+'/ranking');

		function parseTermRanking(termRanking) {
			if(termRanking.Ranking === "" || termRanking.Ranking === undefined) {
			 	return;
			 }

			var newRanking;
			try {
				newRanking = JSON.parse(termRanking.Ranking);
			} catch(e) {
				console.log("Receive invalid JSON ranking form server", e);
				return;
			}

			replaceRanking(newRanking);
		}

		function replaceRanking(r) {
			self.ranking.length = 0;
			Array.prototype.push.apply(self.ranking, r);
		}

		this.stream.onopen = function() {

		};

		this.stream.onmessage = function(e) {
			var bundle = e.data;

			if(bundle.Action == "update") {
				//Often if the user created the item, it will already be in place so treat as an update
				parseTermRanking(bundle.Val);
			}

			$rootScope.$digest();
		};

		this.stream.onclose = function() {

		};

		this.init = function() {
			TaxonomyRankingResource.get(
				{Term: this.term},
				function(termRanking) {
					parseTermRanking(termRanking);
				},
				function(e) {
					console.log(e);
				}
			);
		};

		this.init();

		this.update = function(newRanking, sucess, failure) {
			var oldRanking = this.ranking.slice(0);

			replaceRanking(newRanking);

			TaxonomyRankingResource.update(
				{
					Term: this.term,
					Ranking: JSON.stringify(newRanking)
				},
				function() {
					//Success do nothing!
					if(angular.isFunction(sucess)) sucess();
				},
				function(e) {
					if(e.status != 304) {
						Alerter.error("There was a problem updating the ranking. "+"Status:"+e.status+". Reply Body:"+e.data);
						console.log(e);

						replaceRanking(oldRanking);
					}
					if(angular.isFunction(failure)) failure(e);
				}
			);
		};
	}

	return TaxonomyRankingCache;
})

.factory('TaxonomyResource', function ($resource) {
	var TaxonomyResource = $resource(
		'/1.0/taxonomy/:Term',
		{ Term: '@Term' },
		{
			update: { method: 'PUT', params:{ Term: '@Term' } }
		}
	);

	return TaxonomyResource;
})

.factory('TaxonomyStream', function(SockJS) {
	return new SockJS('/1.0/streaming/taxonomy');
})

.factory('TaxonomyCache', function($rootScope, SockJS, TaxonomyResource, TaxonomyStream) {
	function TaxonomyCache() {
		var self = this;

		this.terms = [];
		this.stream = TaxonomyStream;

		function addUpdateTerm(term) {
			var found = false;

			for(var i = 0;i < self.terms.length;i++) {
				if(self.terms[i].Term == item.Term) {
					for (var key in item) {
						self.terms[i][key] = item[key];
					}
					found = true;
					break;
				}
			}

			if(!found) self.terms.push(item);
		}

		function removeItem(item) {
			for(var i = 0;i < self.terms.length;i++) {
				if(self.terms[i].Term == item.Term) {
					self.terms.remove(i);
					break;
				}
			}
		}

		this.stream.onopen = function() {

		};

		this.stream.onmessage = function(e) {
			var bundle = e.data;

			if(bundle.Action == "create" || bundle.Action == "update") {
				//Often if the user created the item, it will already be in place so treat as an update
				addUpdateTerm(bundle.Val);
			} else if(bundle.Action == "delete") {
				removeItem(bundle.Val);
			} else {
				console.log("Malformed message in Item Stream");
				console.log(e);
			}

			$rootScope.$digest();
		};

		this.stream.onclose = function() {

		};

		this.init = function() {
			ItemResource.query(
				function(newItems) {
					self.terms.length = 0;
					for(var i = 0;i < newItems.length;i++) {
						addUpdateTerm(newItems[i]);
					}
				},
				function(e) {
					console.log("Couldn't load items");
					console.log(e);
				}
			);
		};

		this.init();

		this.create = function(item, sucess, failure) {
			ItemResource.save(
				item,
				function(data) {
					//Do not add the item here since it has no ID, it will be added by the websocket callback
					Alerter.success("Item Created", 2000);
					if(angular.isFunction(sucess)) sucess();
				},
				function(e) {
					Alerter.error("There was a problem creating the item. "+"Status:"+e.status+". Reply Body:"+e.data);
					console.log(e);
					if(angular.isFunction(failure)) failure(e);
				}
			);
		};

		this.update = function(item, sucess, failure) {
			//If the item instanct is already in the array, I assume that you will externally handle 
			//rollback if the update fails. If the item instance is not in the array, this function
			//provides rollback if the update fails
			var handleRollback = this.items.indexOf(item) == -1 ? true : false;
			var present = false;
			var oldItem;

			if(handleRollback) {
				for(var i = 0;i < this.items.length;i++) {
					if(this.items[i].Term == item.Term) {
						oldItem = this.items[i];
						this.items[i] = item;
						present = true;
						break;
					}
				}
				if(!present) { //Odd, I guess we'll add as a new item
					this.items.push(item);
				}
			}

			ItemResource.update(
				item,
				function() {
					Alerter.success("Item Updated", 2000);
					if(angular.isFunction(sucess)) sucess();
					//TODO: Make a revert mechanism here?
				},
				function(e) {
					if(e.status == 304) {
						Alerter.warn("Nothing was changed.", 2000);
					} else {
						Alerter.error("There was a problem updating the item. "+"Status:"+e.status+". Reply Body:"+e.data);
						if(handleRollback && present) {
							for(var i = 0;i < self.terms.length;i++) {
								if(self.terms[i].Term == oldItem.Term) {
									self.terms[i] = oldItem;
									break;
								}
							}
						}
						console.log(e);
					}
					if(angular.isFunction(failure)) failure(e);
				}
			);
		};

		this.remove = function(item, sucess, failure) {
			ItemResource.remove(
				item,
				function() {
					Alerter.success("Item Deleted", 2000);
					removeItem(item);
					if(angular.isFunction(sucess)) sucess();
				},
				function(e) {
					Alerter.error("There was a problem deleting the item. "+"Status:"+e.status+". Reply Body:"+e.data);
					console.log(e);
					if(angular.isFunction(failure)) failure();
				}
			);
		};
	}

	return new ItemCache();
})

.controller('TagSelectorCtrl', function($scope, TaxonomyResource) {
	if($scope.terms !== undefined) {
		var concat = "";
		for(var i = 0;i < $scope.terms.length;i++) {
			concat += $scope.terms[i].Term+",";
		}
		$scope.terms = concat.substring(0, concat.length-1);
	}

	$scope.tagConfig = {
		data: TaxonomyResource.query(),
		multiple: true,
		id: function(item) {
			return item.Term;
		},
		formatResult: function(item) {
			return item.Term;
		},
		formatSelection: function(item) {
			return item.Term;
		},
		createSearchChoice: function(term) {
			return {Term: term};
		},
		matcher: function(term, text, option) {
			return option.Term.toUpperCase().indexOf(term.toUpperCase())>=0;
		},
		initSelection: function (element, callback) {
			var data = [];
			$(element.val().split(",")).each(function (v, a) {
				data.push({Term: a});
			});
			callback(data);
		}
	};
})

.directive('tagSelector', function() {
	return {
		restrict: 'E',
		scope: {
			terms: "="
		},
		templateUrl: '/template/taxonomy/tag-selector.html',
		controller: 'TagSelectorCtrl',
		link: function(scope, element, attrs) {

		}
	};
});