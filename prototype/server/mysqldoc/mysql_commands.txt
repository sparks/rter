MySQL random sample commands
============================

All of these were actually used for our rtER MySQL database.
Note that MySQL is case-insensitive, but keywords are usually capitalized for sanity/clarity.

- Creating a table:

	mysql> CREATE TABLE content (uid INT(64) NOT NULL AUTO_INCREMENT PRIMARY KEY, phone_id VARCHAR(64) NOT NULL, timestamp TIMESTAMP NOT NULL, filepath VARCHAR(256) NOT NULL);

- Adding a column to a table:

	mysql> ALTER TABLE content ADD geolat DECIMAL(9,6);

- Altering column extra properties, in this case removing the timestamp's auto-update by simply not listing it:

	mysql> ALTER TABLE content CHANGE `timestamp` `timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;