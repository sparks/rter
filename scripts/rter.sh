#!/bin/bash
# chkconfig: 2345 20 80
# description: Description comes here....

# Source function library.
. /etc/init.d/functions

start() {
	export RTER_TEMPLATE_DIR='/usr/share/rter/templates/'
	start-stop-daemon --start --exec /usr/share/rter/rter --make-pidfile --pidfile /var/run/rter.pid
}

stop() {
	start-stop-daemon --stop --exec /usr/share/rter/rter --make-pidfile --pidfile /var/run/rter.pid
}

case "$1" in 
	start)
		start
	;;
	
	stop)
		stop
	;;

	retart)
		stop
		start
	;;

	*)
		echo "Usage: $0 {start|stop|restart}"
esac

exit 0
