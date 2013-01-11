MySQL Setup
===========

Run these commands to setup the MySQL databases for the rter project

- - -

    CREATE TABLE content (
         uid INT(64) NOT NULL AUTO_INCREMENT PRIMARY KEY, 
         phone_id VARCHAR(64) NOT NULL, 
         timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
         filepath VARCHAR(256) NOT NULL,
         geolat DECIMAL(9,6),
         geolong DECIMAL(9,6)
    );
    
    CREATE TABLE whitelist (
         uid INT(64) NOT NULL AUTO_INCREMENT PRIMARY KEY, 
         phone_id VARCHAR(64) NOT NULL UNIQUE
    );
    
    CREATE TABLE layout (
         uid INT(64) NOT NULL AUTO_INCREMENT PRIMARY KEY, 
         phone_id VARCHAR(64) NOT NULL UNIQUE,
         col INT(32),
         row INT(32),
         size_x INT(32),
         size_y INT(32)
    );
    
    INSERT INTO whitelist (phone_id) VALUES
        ("fe7f033bfc7b3625fa06c9a3b6b54b2c81eeff98"),
        ("b6200c5cc15cfbddde2874c40952a7aa25a869dd"),
        ("852decd1fbc083cf6853e46feebb08622d653602"),
        ("e1830fcefc3f47647ffa08350348d7e34b142b0b"),
        ("48ad32292ff86b4148e0f754c2b9b55efad32d1e"),
        ("acb519f53a55d9dea06efbcc804eda79d305282e")
    ;