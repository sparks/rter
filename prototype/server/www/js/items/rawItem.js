angular.module('rawItem', ['ui.directives', 'taxonomy'])

.controller('FormRawItemCtrl', function($scope, Taxonomy) {
	//This is kinda terrible
	if($scope.item.Terms !== undefined) {
		var concat = "";
		for(var i = 0;i < $scope.item.Terms.length;i++) {
			concat += $scope.item.Terms[i].Term+",";
		}
		$scope.item.Terms = concat.substring(0, concat.length-1);
	}

	$scope.tagConfig = {
		data: Taxonomy.query(),
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
			$(element.val().split(",")).each(function () {
				data.push({Term: this});
			});
			callback(data);
		}
	};
})

.directive('formRawItem', function() {
	return {
		restrict: 'E',
		scope: {
			item: "=",
			form: "="
		},
		templateUrl: '/template/items/raw/form-raw-item.html',
		controller: 'FormRawItemCtrl',
		link: function(scope, element, attr) {

		}
	};
});
