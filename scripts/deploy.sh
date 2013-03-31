#!/bin/bash
sudo /etc/init.d/rter.sh stop
cd ../prototype/server/
go install rter
cp ./bin/rter /usr/local/share/rter/
cp -r www/ /usr/local/share/rter/
cp -r uploads/ /usr/local/share/rter/
chgrp -R rter /usr/local/share/rter/
chmod g+rw -R  /usr/local/share/rter/
sudo /etc/init.d/rter.sh start
