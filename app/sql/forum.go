package slc

var (
	GetForumBySlug = `SELECT title, "user", slug, posts, threads
					  FROM forums
					  WHERE slug = $1;`

	CreateForum = `INSERT INTO forums
				   (title, "user", slug)
				   VALUES ($1, $2, $3);`

	CheckForumExistBySlug = `SELECT slug
							 FROM forums
							 WHERE slug = $1;`

	ForumPostsCoutnUpdate = `UPDATE forums
							 SET posts = posts + $2
							 WHERE slug = $1;`
)
