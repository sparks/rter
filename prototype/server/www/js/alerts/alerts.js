angular.module('alerts', ['ui.bootstrap', 'ng'])

.factory('Alerter', function ($timeout) {
	function Alerter() {
		this.alerts = [];
		var self = this;

		this.warn = function(msg, timeout) {
			var alert = {msg: msg};
			this.alerts.push(alert);
			if(timeout !== undefined) {
				$timeout(function() {
					self.alerts.remove(self.alerts.indexOf(alert));
				}, timeout);
			}
		};

		this.error = function(msg, timeout) {
			var alert = {type: 'error', msg: msg};
			this.alerts.push(alert);
			if(timeout !== undefined) {
				$timeout(function() {
					self.alerts.remove(self.alerts.indexOf(alert));
				}, timeout);
			}
		};

		this.success = function(msg, timeout) {
			var alert = {type: 'success', msg: msg};
			this.alerts.push(alert);
			if(timeout !== undefined) {
				$timeout(function() {
					self.alerts.remove(self.alerts.indexOf(alert));
				}, timeout);
			}
		};

		this.alert = function(alert, timeout) {
			this.alerts.push(alert);
			if(timeout !== undefined) {
				$timeout(function() {
					self.alerts.remove(self.alerts.indexOf(alert));
				}, timeout);
			}
		};

	}

	return new Alerter();
})

.controller('AlertsCtrl', function($scope, Alerter) {
	$scope.alerts = Alerter.alerts;
	$scope.closeAlert = function(index) {
		$scope.alerts.splice(index, 1);
	};
});