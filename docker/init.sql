use bank_transactions;

CREATE TABLE IF NOT EXISTS tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    access_token  VARCHAR(600) NOT NULL,
	token_type    CHAR(10)      NOT NULL,
	expires_in    INT      NOT NULL,
    when_expires    DATETIME      NOT NULL,
	scope         VARCHAR(1500)     NOT NULL,
	id_token       VARCHAR(600)     NOT NULL,
	refresh_token VARCHAR(600) NOT NULL ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP 
)  ENGINE=INNODB;

INSERT INTO tokens(access_token, token_type, expires_in, when_expires, scope, id_token, refresh_token) VALUES('','',0,'1999-01-01 00:00','','','');

CREATE TABLE IF NOT EXISTS log (
   id INT AUTO_INCREMENT PRIMARY KEY,
   message VARCHAR(200),
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
   updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
 