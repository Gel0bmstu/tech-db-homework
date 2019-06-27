package slc

var (
	CreateThread = `INSERT INTO threads
					(title, author, forum, message, slug, created)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING id;`

	GetThreadBySlug = `SELECT id, title, author, forum, message, votes, slug, "created"
					   FROM threads
					   WHERE slug = $1;`

	GetThreadById = `SELECT id, title, author, forum, message, votes, slug, "created"
					 FROM threads
					 WHERE id = $1;`

	GetThreadsByForumSinceDESC = `SELECT id, title, forum, author, message, slug, created, votes
							  	  FROM threads
							  	  WHERE forum = $1 AND created <= $2 
							  	  ORDER BY created DESC 
							  	  LIMIT $3::TEXT::INTEGER;`

	GetThreadsByForumSince = `SELECT id, title, forum, author, message, slug, created, votes
							  FROM threads
							  WHERE forum = $1 AND created >= $2 
							  ORDER BY created
							  LIMIT $3::TEXT::INTEGER;`

	GetThreadsByForumDESC = `SELECT id, title, forum, author, message, slug, created, votes
							 FROM threads
							 WHERE forum = $1
							 ORDER BY created DESC 
							 LIMIT $2::TEXT::INTEGER;`

	GetThreadsByForum = `SELECT id, title, forum, author, message, slug, created, votes
						 FROM threads
						 WHERE forum = $1
						 ORDER BY created
						 LIMIT $2::TEXT::INTEGER;`

	CheckThreadBySlug = `SELECT id
						 FROM threads
						 WHERE slug = $1;`

	CheckThreadByID = `SELECT id
					   FROM threads
					   WHERE id = $1;`

	GetThreadSlugAndIdBySlug = `SELECT forum, id
								FROM threads
								WHERE slug = $1;`

	GetThreadSlugAndIdByID = `SELECT forum, id
							  FROM threads
							  WHERE id = $1;`

	ThreadUpdate = `UPDATE threads 
					SET message = COALESCE(NULLIF($2, ''), message), title = COALESCE(NULLIF($3, ''), title)
					WHERE id = $1
					RETURNING forum, author, slug, created, message, title, votes;`
)
