angular.module('moment', [])

.factory('moment', function () {
	return moment;
})

.filter('fromNow', function(moment) {
	return function(input) {
		return moment(input).fromNow();
	};
});