package slc

var (
	CreateIndexOnPerf = `
		-- Useer:
		CREATE INDEX IF NOT EXISTS UserNicknameSelectIdx ON users (nickname);
		
		-- Forums:
		CREATE INDEX IF NOT EXISTS ForumSlugAndIdIdx ON forums (slug, id);

		-- Threads:
		CREATE INDEX IF NOT EXISTS ThreadCreatedAndForumIdx ON threads (created, forum);

		-- Posts
		CREATE INDEX IF NOT EXISTS PostsThreadAndIdIdx ON posts (thread, id);
		CREATE INDEX IF NOT EXISTS PostParentsThreadAndIdIdx ON posts (thread, id) WHERE parent = 0;
		CREATE INDEX IF NOT EXISTS PostParentPathIdx ON posts (thread, (path || array[id]), parent);
		
		-- Cluster tables
		CLUSTER users USING UserNicknameSelectIdx;
		CLUSTER forums USING ForumSlugAndIdIdx;
		CLUSTER threads USING ThreadCreatedAndForumIdx;
		-- CLUSTER posts USING PostParentsThreadAndIdIdx;

	`
)
