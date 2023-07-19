package myPac

type MigrationQuery struct {
	Version int
	Queries []string
	Message string
}

var MigrationQueries []MigrationQuery = []MigrationQuery{
	{
		Version: 1,
		Message: "creating tables: tbuser - tbmovie - tb_config",
		Queries: []string{`
			CREATE TABLE tbuser (
				username varchar NOT NULL,
				passhash varchar NOT NULL,
				CONSTRAINT tbuser_pk PRIMARY KEY (username)
			);`,
			/* create new user for test with pass 123456 */ `
			INSERT INTO tbuser
				(username, passhash)
			VALUES
				('sadeq', 'MTIzNDU247DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=');`, `
			CREATE TABLE tbmovie (
				username varchar NOT NULL,
				movie_name varchar NOT NULL,
				"year" varchar NOT NULL,
				CONSTRAINT tbmovie_pk PRIMARY KEY (username, movie_name, year),
				CONSTRAINT tbmovie_fk FOREIGN KEY (username) REFERENCES tbuser(username) ON DELETE CASCADE ON UPDATE CASCADE
			);`, `
			CREATE TABLE public.tb_config (
				migration_version int NOT NULL DEFAULT 0,
				last_migration timestamp NOT NULL
			);`, `
			INSERT INTO tb_config (migration_version, last_migration) VALUES(1, now());`,
		},
	}, {
		Version: 2,
		Message: "add email",
		Queries: []string{`
			ALTER TABLE public.tbuser ADD email varchar NULL;`, `
			ALTER TABLE public.tbuser ADD email_verified boolean NOT NULL DEFAULT false;`, `
			UPDATE tbuser SET email=concat(tbuser.username, '@no-mail.com'), email_verified=true;`, `
			ALTER TABLE public.tbuser ALTER COLUMN email SET NOT NULL;`, `
			ALTER TABLE public.tbuser ADD CONSTRAINT tbuser_un_email UNIQUE (email);`, `
			UPDATE public.tbuser SET email='sadeq@fav-mov.com' , email_verified=true WHERE username='sadeq';`, `
			UPDATE tb_config SET migration_version=2, last_migration=now();`,
		},
	}, {
		Version: 3,
		Message: "test",
		Queries: []string{`
			UPDATE tb_config SET migration_version=3, last_migration=now();`,
		},
	},
}
