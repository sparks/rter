# Quick and Dirty Server Setup

 * Install go and Mysql
 * Create a database called 'rter'
 	* Log in to MYSQL as root mysql user and run this commands:
 	* `CREATE DATABASE rter;`
 * You can delete existing databases if you need
 	* Log in to MYSQL as root mysql user and run this commands:
 	* `DROP DATABASE nameofdb;`
 * Create new mysql user
 	* Log in to MYSQL as root mysql user and run these commands:
 	* `CREATE USER 'rter'@'localhost' IDENTIFIED BY 'j2pREch8';`
 	* `GRANT ALL PRIVILEGES ON rter . * TO 'rter'@'localhost';`
 	* `FLUSH PRIVILEGES;`
 * Setup Databases by running this command from the projects 'scripts' directory
	* `mysql -u rter -pj2pREch8 rter < ./rter_v2_datefix.sql`
 * Add the appropriate directories to your GOPATH, should look vaguely like this in your ~/.bash_rc or ~/.bash_profile
 	* `export GOPATH='/Path/to/rter/prototype/server':$GOPATH`
	* `export GOPATH='/Path/to/rter/prototype/videoserver':$GOPATH`
 * Go into the server direction 'prototypes/server' and launch the server
 	* `go run src/rter/rter.go`
 * For more options you can type 
	* `go run src/rter/rter.go --help`
 * You can install missing go libraries by typing
 	* `go install name/of/lib`
