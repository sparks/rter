#!/bin/bash
sudo /etc/init.d/videoserver.sh stop
cd ../prototype/videoserver/
go install videoserver 
cp ./bin/videoserver /usr/local/share/videoserver/
chgrp -R rter /usr/local/share/videoserver/
chmod g+rw -R  /usr/local/share/videoserver/
sudo /etc/init.d/videoserver.sh start
