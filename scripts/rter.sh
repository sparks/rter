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
	export RTER_DIR='/usr/local/share/rter/'
	export RTER_LOGFILE='/var/log/rter.log'
	setcap 'cap_net_bind_service=+ep' /usr/local/share/rter/rter
	start-stop-daemon --start --background --make-pidfile --pidfile /var/run/rter.pid --chuid rter --exec /usr/local/share/rter/rter -- --http-port 80 $*
}

stop() {
	start-stop-daemon --stop --pidfile /var/run/rter.pid --exec /usr/local/share/rter/rter
}

case "$1" in 
	start)
		shift
		start $*
	;;
	
	stop)
		stop
	;;

	restart)
		stop
		shift
		start $*
	;;

	*)
		echo "Usage: $0 {start|stop|restart}"
esac

exit 0
