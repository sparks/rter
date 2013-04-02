#!/bin/bash
#/etc/init.d/videserver.sh
#
### BEGIN INIT INFO
# Provides:          videoserver
# Required-Start:    $all
# Required-Stop:     $all
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Run the videoserver
# Description:       Run the videoserver developped for rtER
### END INIT INFO

start() {
        start-stop-daemon --start --background --make-pidfile --pidfile /var/run/videoserver.pid --chuid rter --exec /usr/local/share/videoserver/videoserver -- --config /usr/local/share/videoserver/config.json
}

stop() {
        start-stop-daemon --stop --pidfile /var/run/videoserver.pid --exec /usr/local/share/videoserver/videoserver
}

case "$1" in
        start)
                start
        ;;

        stop)
                stop
        ;;

        restart)
                stop
                start
        ;;

        *)
                echo "Usage: $0 {start|stop|restart}"
esac

exit 0

