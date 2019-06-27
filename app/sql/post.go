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
								 WHERE thread = $1 AND id < $2::TEXT::INTEGER
								 ORDER BY id DESC
								 LIMIT $3::TEXT::INTEGER;`

	GetPostsByIdFlatSince = `SELECT id, author, message, thread, created, parent
							 FROM posts
							 WHERE thread = $1 AND id > $2::TEXT::INTEGER
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
									 WHERE id = $2::TEXT::INTEGER
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
								 WHERE id = $2::TEXT::INTEGER
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
	GetPostsByIdParentTreeSinceDesc = `SELECT p.id, p.author, p.message, p.thread, p.created, p.parent 
									   FROM posts p
									   JOIN (
										   SELECT p1.id 
										   FROM posts p1
										   WHERE p1.thread = $1 AND p1.parent = 0 AND p1.id < (
											   SELECT p2.path[1] 
											   FROM posts p2 
											   WHERE p2.id = $2)  
										   ORDER BY p1.id DESC 
										   LIMIT $3::TEXT::INTEGER) 
									   AS rootParents 
									   ON (
										   rootParents.id = p.path[1])
									   ORDER BY p.path[1] DESC, p.path;`

	GetPostsByIdParentTreeSince = `SELECT p.id, p.author, p.message, p.thread, p.created, p.parent 
								   FROM posts p
								   JOIN (
									   SELECT id
									   FROM posts p1
									   WHERE p1.thread = $1 AND p1.parent = 0 AND p1.id > (
										   SELECT p2.path[1] 
										   FROM posts p2
										   WHERE p2.id = $2)  
									   ORDER BY id 
									   LIMIT $3::TEXT::INTEGER) 
								   AS rootParents 
								   ON (
									   rootParents.id = p.path[1])
								   ORDER BY p.path[1], p.path;`

	GetPostsByIdParentTreeDesc = `SELECT p.id, p.author, p.message, p.thread, p.created, p.parent 
								  FROM posts p
								  JOIN (
									  SELECT id AS rootParentsId 
									  FROM posts 
									  WHERE thread = $1 AND parent = 0
									  ORDER BY id DESC 
									  LIMIT $2::TEXT::INTEGER)
								  AS rootParents 
								  ON (
									  rootParents.rootParentsId = p.path[1])
								  ORDER BY path[1] DESC, path;`

	GetPostsByIdParentTree = `SELECT p.id, p.author, p.message, p.thread, p.created, p.parent 
							  FROM posts p
							  JOIN (
								  SELECT id as rootParentsId
								  FROM posts 
								  WHERE thread = $1 AND parent = 0
								  ORDER BY id
								  LIMIT $2::TEXT::INTEGER) 
							  AS rootParents 
							  ON (
								  rootParents.rootParentsId = p.path[1]) 
							  ORDER BY path[1], path;`
)
