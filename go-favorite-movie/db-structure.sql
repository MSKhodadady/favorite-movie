CREATE TABLE tbuser (
	username varchar NOT NULL,
	passhash varchar NOT NULL,
	CONSTRAINT tbuser_pk PRIMARY KEY (username)
);

CREATE TABLE tbmovie (
	username varchar NOT NULL,
	movie_name varchar NOT NULL,
	"year" varchar NOT NULL,
	CONSTRAINT tbmovie_pk PRIMARY KEY (username, movie_name, year),
	CONSTRAINT tbmovie_fk FOREIGN KEY (username) REFERENCES tbuser(username) ON DELETE CASCADE ON UPDATE CASCADE
);