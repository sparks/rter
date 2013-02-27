-- MySQL rter v1
-- ===========
-- Run these commands to setup the MySQL databases for the rter v1 project

DROP TABLE IF EXISTS content;
CREATE TABLE IF NOT EXISTS content (
	uid INT(64) NOT NULL AUTO_INCREMENT,
	content_id VARCHAR(64) NOT NULL,
	content_type VARCHAR(64) NOT NULL,
	timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	filepath VARCHAR(256) NOT NULL DEFAULT "",
	geolat DECIMAL(9,6) NOT NULL DEFAULT 0,
	geolng DECIMAL(9,6) NOT NULL DEFAULT 0,
	heading DECIMAL(9, 6) NOT NULL DEFAULT 0,
	description TEXT NOT NULL DEFAULT "",
	url VARCHAR(256) NOT NULL DEFAULT "",
	PRIMARY KEY(uid),
	KEY (content_id, timestamp)
);

DROP TABLE IF EXISTS phones;
CREATE TABLE IF NOT EXISTS phones (
	uid INT(64) NOT NULL AUTO_INCREMENT,
	phone_id VARCHAR(64) NOT NULL UNIQUE,
	target_heading DECIMAL(9, 6) NOT NULL DEFAULT 0,
	PRIMARY KEY(uid)
);

DROP TABLE IF EXISTS layout;
CREATE TABLE IF NOT EXISTS layout (
	uid INT(64) NOT NULL AUTO_INCREMENT,
	content_id VARCHAR(64) NOT NULL UNIQUE KEY,
	col INT(32) NOT NULL,
	row INT(32) NOT NULL,
	size_x INT(32) NOT NULL DEFAULT 1,
	size_y INT(32) NOT NULL DEFAULT 1,
	PRIMARY KEY(uid)
);

INSERT INTO phones (phone_id) VALUES
	("1e7f033bfc7b3625fa07c9a3b6b54d2c81eeff98"),
	("fe7f033bfc7b3625fa06c9a3b6b54b2c81eeff98"),
	("b6200c5cc15cfbddde2874c40952a7aa25a869dd"),
	("852decd1fbc083cf6853e46feebb08622d653602"),
	("e1830fcefc3f47647ffa08350348d7e34b142b0b"),
	("48ad32292ff86b4148e0f754c2b9b55efad32d1e"),
	("acb519f53a55d9dea06efbcc804eda79d305282e"),
	("ze7f033bfc7b3625fa06c5a316b54b2c81eeff98"),
	("t6200c5cc15cfbddde2875c41952a7aa25a869dd"),
	("952decd1fbc083cf6853e56f1ebb08622d653602"),
	("y1830fcefc3f47647ffa05351348d7e34b142b0b"),
	("x8ad32292ff86b4148e0f55412b9b55efad32d1e"),
	("qcb519f53a55d9dea06ef5cc104eda79d305282e")
;

-- delete from content where uid >= 0;alter table content AUTO_INCREMENT=1;