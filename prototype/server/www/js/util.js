Array.prototype.remove = function(from, to) {
	var rest = this.slice((to || from) + 1 || this.length);
	this.length = from < 0 ? this.length + from : from;
	return this.push.apply(this, rest);
};

function map(val, start_min, start_max, end_min, end_max) {
	return (val-start_min)/(start_max-start_min)*(end_max-end_min)+end_min;
}