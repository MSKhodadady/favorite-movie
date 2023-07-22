package myPac

const (
	QFirst20UserName = `
	select 
		username, count(movie_name) movie_count
	from (
		select
			t.username , t2.movie_name 
		from tbuser t natural left outer join tbmovie t2-- on t.username  = t2.username
	) q
	group by username
	order by movie_count desc
	limit 20;
	`

	QSignIn = `SELECT username, passHash 
	FROM tbUser WHERE username = $1`

	QAddUser = `INSERT INTO tbUser (username, passhash, email)
	VALUES ($1, $2, $3);`

	QAddMovie = `INSERT INTO public.tbmovie (username, movie_name, "year")
	VALUES($1, $2, $3);`

	QDeleteMovie = `DELETE FROM tbMovie
	WHERE username=$1 AND movie_name=$2 AND "year"=$3;`
)
