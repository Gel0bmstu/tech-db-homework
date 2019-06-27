package slc

var (
	GetUserByNicknameOrEmail = `SELECT nickname, fullname, about, email
							 	FROM users
								WHERE nickname = $1 OR email = $2;`

	GetUserByNickname = `SELECT nickname, fullname, about, email
						 FROM users
						 WHERE nickname = $1;`

	GetUserByEmail = `SELECT nickname, fullname, about, email
					  FROM users
					  WHERE email = $1;`

	CreateUser = `INSERT INTO users (nickname, fullname, about, email)
				  VALUES ($1, $2, $3, $4);`

	UserHelpTableInsert = `INSERT INTO user_forum (nickname, fullname, about, email, forum)
						   VALUES ($1, $2, $3, $4, $5)
						   ON CONFLICT DO NOTHING;`

	UpdateUser = `UPDATE users 
				  SET fullname = COALESCE(NULLIF($2, ''), fullname), about = COALESCE(NULLIF($3, ''), about), email = COALESCE(NULLIF($4, ''), email) 
				  WHERE nickname = $1 
				  RETURNING fullname, about, email;`

	CheckUserNicknameByNickname = `SELECT nickname
								   FROM users
								   WHERE nickname = $1;`

	CheckUserEmailByEmail = `SELECT email
							 FROM users
							 WHERE email = $1;`

	// GetForumUsersSinceDesc = `SELECT u.nickname, u.fullname, u.about, u.email
	// 						  FROM user_forum uf
	// 						  JOIN users u ON (
	// 							  u.nickname = uf.nickname)
	// 						  WHERE uf.forum = $1 AND u.nickname < $2
	// 						  ORDER BY u.nickname DESC
	// 						  `

	// GetForumUsersSince = `SELECT u.nickname, u.fullname, u.about, u.email
	// 					  FROM user_forum uf
	// 					  JOIN users u ON (
	// 						  u.nickname = uf.nickname)
	// 					  WHERE  uf.forum = $1 AND u.nickname > $2
	// 					  ORDER BY u.nickname
	// 					  `

	// GetForumUsersDesc = `SELECT u.nickname, u.fullname, u.about, u.email
	// 					 FROM user_forum uf
	// 					 JOIN users u ON (
	// 						 u.nickname = uf.nickname)
	// 					 WHERE uf.forum = $1
	// 					 ORDER BY u.nickname DESC
	// 					 `

	// GetForumUsers = `SELECT u.nickname, u.fullname, u.about, u.email
	// 				 FROM user_forum uf
	// 				 JOIN users u ON (
	// 					 u.nickname = uf.nickname)
	// 				 WHERE uf.forum = $1
	// 				 ORDER BY u.nickname
	// 				 `

	GetForumUsersSinceDesc = `SELECT nickname, fullname, about, email
							  FROM user_forum
							  WHERE forum = $1 AND nickname < $2
							  ORDER BY nickname DESC 
							  `

	GetForumUsersSince = `SELECT nickname, fullname, about, email
						  FROM user_forum
						  WHERE  forum = $1 AND nickname > $2
						  ORDER BY nickname
						  `

	GetForumUsersDesc = `SELECT nickname, fullname, about, email
						 FROM user_forum
						 WHERE forum = $1
						 ORDER BY nickname DESC 
						 `

	GetForumUsers = `SELECT nickname, fullname, about, email
					 FROM user_forum
					 WHERE forum = $1
					 ORDER BY nickname
					 `
)
