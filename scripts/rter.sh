#!/bin/bash
#/etc/init.d/rter.sh
#
### BEGIN INIT INFO
# Provides:          rter
# Required-Start:    $all
# Required-Stop:     $all
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Run the rtER web server 
# Description:       Run the rtER web server for client/mobile 
### END INIT INFO

start() {
	export RTER_DIR='/usr/share/rter/'
	start-stop-daemon --start --background --exec /usr/share/rter/rter --make-pidfile --pidfile /var/run/rter.pid
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
