package slc

var (
	GetService = `SELECT (
				  SELECT count(*) 
				  FROM forums) as forums,  
				  (SELECT count(*)
				  FROM posts) as posts,
				  (SELECT count(*)
				  FROM users) as users,
				  (SELECT count(*)
				  FROM threads) as threads;`

	ClearDb = `TRUNCATE votes, posts, threads, forums, users 
			   RESTART IDENTITY CASCADE;`
)
