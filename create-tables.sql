DROP TABLE IF EXISTS heroi;
CREATE TABLE heroi (
  id         INT AUTO_INCREMENT NOT NULL,
  nome       VARCHAR(255) NOT NULL,
  classe     CHAR(1) NOT NULL,
  ranking    SMALLINT(2) NOT NULL,
  PRIMARY KEY (`id`)
);
