#!/bin/bash
sudo /etc/init.d/rter.sh stop
cd ../prototype/server/
go build -o rter main.go
mv ./rter /usr/local/share/rter/
cp -r templates/ /usr/local/share/rter/
chgrp -R rter /usr/local/share/rter/
sudo /etc/init.d/rter.sh start
