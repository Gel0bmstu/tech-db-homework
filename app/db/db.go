package database

import (
	"github.com/jackc/pgx"
)

var (
	DB        *pgx.ConnPool
	Timestamp = "2006-01-02T18:04:05.01+03:00"
)

func ResetDB() (err error) {
	sql := `
		DROP TABLE IF EXISTS "forums" CASCADE;
		DROP TABLE IF EXISTS "users" CASCADE;
		DROP TABLE IF EXISTS "threads" CASCADE;
		DROP TABLE IF EXISTS "posts" CASCADE;
		DROP TABLE IF EXISTS "votes" CASCADE;
		DROP TABLE IF EXISTS "user_forum" CASCADE;

		-- ---------------------------------------------------------------------
		-- TABLES:

		CREATE UNLOGGED TABLE "users" (
			"id" BIGSERIAL PRIMARY KEY,
			"nickname" citext UNIQUE COLLATE ucs_basic,	
			"fullname" text NOT NULL,
			"about" text,
			"email" citext UNIQUE
		);

		CREATE UNLOGGED TABLE "forums"(
			"id" BIGSERIAL PRIMARY KEY,
			"posts" integer DEFAULT 0,
			"user" citext REFERENCES users(nickname), 
			"slug" citext NOT NULL UNIQUE,
			"threads" integer DEFAULT 0,
			"title" text
		);
		CREATE UNLOGGED TABLE "threads"(
			"id" BIGSERIAL PRIMARY KEY,
			"author" citext REFERENCES users(nickname),
			"created" timestamp WITH TIME ZONE DEFAULT NOW(),
			"forum" citext NOT NULL REFERENCES forums(slug),
			"message" text NOT NULL,
			"slug" citext UNIQUE,
			"title" text,
			"votes" integer DEFAULT 0 NOT NULL
		);

		CREATE UNLOGGED TABLE "posts" (
			"id" BIGSERIAL,
			"parent" integer DEFAULT 0,
			"author" citext,
			"message" text NOT NULL,			
			"isedited" boolean DEFAULT 'false',
			"forum" citext REFERENCES forums(slug),
			"thread" bigint REFERENCES threads(id),
			-- "created" timestamp WITH TIME ZONE DEFAULT '1970-01-01T00:00:00.000Z',
			"path" bigint[],
			"childs" bigint[]
		);

		CREATE UNLOGGED TABLE "votes" (
			"id" BIGSERIAL PRIMARY KEY,
			UNIQUE ("nickname", thread),
			"nickname" citext NOT NULL REFERENCES users(nickname),
			"thread" bigint REFERENCES threads(id),
			"voice" int2 
		);

		CREATE UNLOGGED TABLE "user_forum" (
			"nickname" citext COLLATE ucs_basic REFERENCES users(nickname),
			"forum" citext REFERENCES forums(slug),
			"fullname" text NOT NULL,
			"about" text,
			"email" citext,
			CONSTRAINT user_forum_const UNIQUE ("nickname", "forum")
		);

		-- ---------------------------------------------------------------------
		-- INDEXES:
		
		-- users:
		
		CREATE INDEX IF NOT EXISTS idxUserAllSelect 
			ON users (id, nickname, fullname, about, email);	
			
		CREATE INDEX IF NOT EXISTS idxUserForumAll 
			ON user_forum (nickname, forum);

		-- forums:
		
		CREATE INDEX IF NOT EXISTS idxForumSlugId
			ON forums (slug, id);

		-- threads:
		
		CREATE INDEX IF NOT EXISTS idxThreadSlugId 
			ON threads (slug, id);

		CREATE INDEX IF NOT EXISTS idxThreadForumCreated
			ON threads (forum, created);

		-- posts:

		CREATE INDEX IF NOT EXISTS idxPostIdThread
			ON posts (id, thread);

		CREATE INDEX IF NOT EXISTS idxPostPath 
			ON posts ((path[1]), path);

		CREATE INDEX IF NOT EXISTS idxPostThreadPath
			ON posts (thread, path);
<<<<<<< HEAD

		CREATE INDEX IF NOT EXISTS idxPostRootParentThreadId
			ON posts(thread, id) 
			WHERE parent = 0;

=======
			
		CREATE INDEX IF NOT EXISTS idxPostRootParentThreadId
			ON posts(thread, id) 
			WHERE parent = 0;
>>>>>>> cd9c53f27961de658a698e66e3e6b75b96cb5e28
		-- ---------------------------------------------------------------------
		-- HELP FUNCTIONS:

		-- Threads:

		CREATE OR REPLACE FUNCTION threadInsert() RETURNS TRIGGER AS $thread_insert$
			BEGIN
				UPDATE forums 
				SET threads = threads + 1
				WHERE slug = NEW.forum;

				-- INSERT INTO user_forum (nickname, forum)
				-- VALUES (NEW.Author, NEW.Forum)
				-- ON CONFLICT DO NOTHING;
				RETURN NULL;
			END;
		$thread_insert$ LANGUAGE plpgsql;

		CREATE OR REPLACE FUNCTION voteInsert() RETURNS TRIGGER AS $vote_insert$
			BEGIN
				UPDATE threads 
				SET votes = votes + NEW.voice
				WHERE id = NEW.thread;
				RETURN NULL;
			END;
		$vote_insert$ LANGUAGE plpgsql;

		CREATE OR REPLACE FUNCTION voteUpdate() RETURNS TRIGGER AS $vote_update$
			BEGIN
				UPDATE threads 
				SET votes = votes - OLD.voice + NEW.voice
				WHERE id = NEW.thread;
				RETURN NULL;
			END;
		$vote_update$ LANGUAGE plpgsql;

		CREATE OR REPLACE FUNCTION userForumInsert() RETURNS TRIGGER AS $user_forum_insert$
		BEGIN
			INSERT INTO user_forum (nickname, forum)
			VALUES (NEW.Author, NEW.Forum)
			ON CONFLICT DO NOTHING;
			RETURN NULL;
		END;
		$user_forum_insert$ LANGUAGE plpgsql;

		CREATE OR REPLACE FUNCTION changePath() RETURNS TRIGGER AS $change_path$
			DECLARE
				prevPath bigint[];
			BEGIN
				IF (NEW.parent <> 0) AND (NEW.parent IS NOT NULL) 
					THEN
						SELECT path 
						FROM posts
						WHERE id = NEW.parent
						INTO prevPath;

						NEW.path = NEW.path || prevPath || NEW.id;

						UPDATE posts
						SET childs = childs || NEW.id
						WHERE id = NEW.parent;

						-- UPDATE posts
						-- SET path = path || prevPath || NEW.id;
					ELSE
						-- UPDATE posts 
						-- SETid path = path || NEW.id
						-- WHERE id = NEW.parent;

						NEW.path = NEW.path || NEW.id;
					END IF;
				RETURN NEW;
			END;
		$change_path$ LANGUAGE plpgsql;

		-- ---------------------------------------------------------------------
		-- TRIGGERS:

		DROP TRIGGER IF EXISTS voteCreate ON votes;
		DROP TRIGGER IF EXISTS voteUpdate ON votes;
		DROP TRIGGER IF EXISTS threadCreate ON threads;
		DROP TRIGGER IF EXISTS postChildAdd ON posts;
		DROP TRIGGER IF EXISTS postCreate ON posts;

		-- Votes:

		CREATE TRIGGER voteCreate 
			AFTER INSERT 
			ON votes
			FOR EACH ROW EXECUTE PROCEDURE voteInsert();
		
		CREATE TRIGGER voteUpdate
			AFTER UPDATE 
			ON votes
			FOR EACH ROW EXECUTE PROCEDURE voteUpdate();

		-- Threads:

		CREATE TRIGGER threadCreate
			AFTER INSERT 
			ON threads
			FOR EACH ROW EXECUTE PROCEDURE threadInsert();

		-- Posts:

		CREATE TRIGGER changePath
			BEFORE INSERT
			ON posts
			FOR EACH ROW EXECUTE PROCEDURE changePath();
	`
	// fmt.Println(sql)
	_, err = DB.Exec(sql)

	return
}

func Initialize(dbCfg string) (err error) {

	var pgxCfg pgx.ConnConfig

	pgxCfg, err = pgx.ParseURI(dbCfg)

	if err != nil {
		return
	}

	DB, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgxCfg,
		MaxConnections: 100,
	})

	if err != nil {
		return
	}

	err = ResetDB()

	return
}
