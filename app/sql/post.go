package slc

var (
	CheckExistingPostByThreadId = `SELECT id
							 	   FROM posts
								   WHERE thread = $1 AND id = $2;`

	CheckPostByIdAndGetMessege = `SELECT id, message
								  FROM posts
								  WHERE id = $1;`

	CreatePost = `INSERT INTO posts
				  (author, message, parent, thread, forum)
				  VALUES ($1, $2, $3, $4, $5)
				  RETURNING created, id;`

	GetPostById = `SELECT id, author, message, thread, created, parent, forum, isedited
				   FROM posts
				   WHERE id = $1;`

	UpdatePostById = `UPDATE posts
					  SET message = $2, isedited = true 
					  WHERE id = $1
					  RETURNING id, author, message, thread, created, parent, forum, isedited;`

	// Flat sort
	GetPostsByIdFlatSinceDesc = `SELECT id, author, message, thread, created, parent
								 FROM posts
								 WHERE thread = $1 AND id < $2
								 ORDER BY id DESC
								 LIMIT $3::TEXT::INTEGER;`

	GetPostsByIdFlatSince = `SELECT id, author, message, thread, created, parent
							 FROM posts
							 WHERE thread = $1 AND id > $2
							 ORDER BY id
							 LIMIT $3::TEXT::INTEGER;`

	GetPostsByIdFlatDesc = `SELECT id, author, message, thread, created, parent
							FROM posts
							WHERE thread = $1
							ORDER BY id DESC
							LIMIT $2::TEXT::INTEGER;`

	GetPostsByIdFlat = `SELECT id, author, message, thread, created, parent
						FROM posts
						WHERE thread = $1
						ORDER BY id
						LIMIT $2::TEXT::INTEGER;`

	// Tree sort
	GetPostsByIdTreeSinceDesc = `SELECT id, author, message, thread, created, parent 
								 FROM posts 
								 WHERE thread = $1 AND path < 
								 	(
									 SELECT path
									 FROM posts 
									 WHERE id = $2
									) 
								 ORDER BY path DESC 
								 LIMIT $3::TEXT::INTEGER;`

	GetPostsByIdTreeSince = `SELECT id, author, message, thread, created, parent 
							 FROM posts 
							 WHERE thread = $1 
							 AND path > 
							 	(
								 SELECT path
								 FROM posts 
								 WHERE id = $2
								) 
							 ORDER BY path 
							 LIMIT $3::TEXT::INTEGER;`

	GetPostsByIdTreeDesc = `SELECT id, author, message, thread, created, parent 
							FROM posts 
							WHERE thread = $1 
							ORDER BY path DESC 
							LIMIT $2::TEXT::INTEGER;`

	GetPostsByIdTree = `SELECT id, author, message, thread, created, parent 
						FROM posts 
						WHERE thread = $1 
						ORDER BY path 
						LIMIT $2::TEXT::INTEGER;`

	// ParentTree sort
	GetPostsByIdParentTreeSinceDesc = `SELECT id, author, message, thread, created, parent 
									   FROM posts 
									   JOIN (
										   SELECT id AS rootParentsId 
										   FROM posts 
										   WHERE thread = $1 AND parent = 0 AND path[1] < (
											   SELECT path[1] 
											   FROM posts 
											   WHERE id = $2)  
										   ORDER BY id DESC 
										   LIMIT $3::TEXT::INTEGER) 
									   AS rootParents 
									   ON (
										   thread = $1 AND rootParents.rootParentsId = path[1])
									   ORDER BY rootParents.rootParentsId DESC, path;`

	GetPostsByIdParentTreeSince = `SELECT id, author, message, thread, created, parent 
									   FROM posts 
									   JOIN (
										   SELECT id AS rootParentsId 
										   FROM posts 
										   WHERE thread = $1 AND parent = 0 AND path[1] > (
											   SELECT path[1] 
											   FROM posts 
											   WHEre id = $2)  
										   ORDER BY id DESC 
										   LIMIT $3::TEXT::INTEGER) 
									   AS rootParents 
									   ON (
										   thread = $1 AND rootParents.rootParentsId = path[1])
									   ORDER BY rootParents.rootParentsId, path;`

	GetPostsByIdParentTreeDesc = `SELECT id, author, message, thread, created, parent 
								  FROM posts 
								  JOIN (
									  SELECT id AS rootParentsId 
									  FROM posts 
									  WHERE thread = $1 AND parent = 0
									  ORDER BY id DESC LIMIT $2::TEXT::INTEGER)
								  AS rootParents 
								  ON (
									  thread = $1 AND rootParents.rootParentsId = path[1])
								  ORDER BY rootParents.rootParentsId DESC, path;`

	GetPostsByIdParentTree = `SELECT id, author, message, thread, created, parent 
									FROM posts pst
									JOIN (
										SELECT id as rootParentsId
										FROM posts 
										WHERE thread = $1 AND parent = 0
										ORDER BY id
										LIMIT $2::TEXT::INTEGER) 
									AS rootParents 
									ON (
										thread = $1 AND rootParents.rootParentsId = path[1]) 
									ORDER BY path;`
)
