angular.module('cache', [
	'alerts', //Alerts for item actions
	'ng'      //rootscope	
])

.factory('CacheBuilder', function ($rootScope, Alerter) {
	function CacheBuilder(name, resource, stream, matchFn, uncloseable) {
		var self = this;

		this.name = name;
		
		this.stream = stream;
		this.resource = resource;

		this.matchFn = matchFn;

		this.uncloseable = uncloseable;

		this.contents = [];

		function addUpdateElem(newElem) {
			var found = false;

			for(var i = 0;i < self.contents.length;i++) {
				if(self.matchFn(self.contents[i], newElem)) {
					for (var key in newElem) {
						self.contents[i][key] = newElem[key];
					}
					found = true;
					break;
				}
			}

			if(!found) self.contents.push(newElem);
		}

		function removeElem(elem) {
			for(var i = 0;i < self.contents.length;i++) {
				if(self.matchFn(self.contents[i], elem)) {
					self.contents.remove(i);
					break;
				}
			}
		}

		this.stream.onopen = function() {

		};

		this.stream.onmessage = function(e) {
			var bundle = e.data;

			if(bundle.Action == "create" || bundle.Action == "update") {
				addUpdateElem(bundle.Val);
			} else if(bundle.Action == "delete") {
				removeElem(bundle.Val);
			} else {
				console.log("Malformed message in "+self.name+" Stream");
				console.log(e);
			}

			$rootScope.$digest();
		};

		this.stream.onclose = function() {

		};

		this.close = function() {
			if(this.uncloseable) return;
			this.stream.close();
		}

		this.init = function() {
			self.resource.query(
				function(newElems) {
					self.contents.length = 0;
					for(var i = 0;i < newElems.length;i++) {
						addUpdateElem(newElems[i]);
					}
				},
				function(e) {
					if(e.status != 404) {
						console.log("Couldn't load init for "+self.name);
						console.log(e);
					}
				}
			);
		};

		this.init();

		this.create = function(elem, sucess, failure) {
			self.resource.save(
				elem,
				function(data) {
					Alerter.success(self.name+" Created", 2000);
					if(angular.isFunction(sucess)) sucess();
				},
				function(e) {
					Alerter.error("There was a problem creating the "+self.name+". Status:"+e.status+". Reply Body:"+e.data);
					console.log(e);
					if(angular.isFunction(failure)) failure(e);
				}
			);
		};

		this.update = function(elem, sucess, failure) {
			//If the elem instanct is already in the array, I assume that you will externally handle 
			//rollback if the update fails. If the elem instance is not in the array, this function
			//provides rollback if the update fails
			var handleRollback = this.contents.indexOf(elem) == -1 ? true : false;
			var present = false;
			var oldElem;

			if(handleRollback) {
				for(var i = 0;i < this.contents.length;i++) {
					if(self.matchFn(this.contents[i], elem)) {
						oldElem = this.contents[i];
						this.contents[i] = elem;
						present = true;
						break;
					}
				}
				if(!present) { //Odd, I guess we'll add as a new elem
					this.contents.push(elem);
				}
			}

			self.resource.update(
				elem,
				function() {
					Alerter.success(self.name+" Updated", 2000);
					if(angular.isFunction(sucess)) sucess();
					//TODO: Make a revert mechanism here?
				},
				function(e) {
					if(e.status == 304) {
						Alerter.warn("Nothing was changed.", 2000);
					} else {
						Alerter.error("There was a problem updating the "+self.name+". "+"Status:"+e.status+". Reply Body:"+e.data);
						if(handleRollback && present) {
							for(var i = 0;i < self.contents.length;i++) {
								if(self.matchFn(self.contents[i], oldElem)) {
									self.contents[i] = oldElem;
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

		this.remove = function(elem, sucess, failure) {
			self.resource.remove(
				elem,
				function() {
					Alerter.success(self.name+" Deleted", 2000);
					removeElem(elem);
					if(angular.isFunction(sucess)) sucess();
				},
				function(e) {
					Alerter.error("There was a problem deleting the "+self.elem+". "+"Status:"+e.status+". Reply Body:"+e.data);
					console.log(e);
					if(angular.isFunction(failure)) failure();
				}
			);
		};
	}

	return CacheBuilder;
})
