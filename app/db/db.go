package database

import (
	"fmt"

	"github.com/jackc/pgx"
)

var (
	DB *pgx.ConnPool
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
		) WITH (autovacuum_enabled=false);

		CREATE UNLOGGED TABLE "forums"(
			"id" BIGSERIAL PRIMARY KEY,
			"posts" integer DEFAULT 0,
			"user" citext REFERENCES users(nickname), 
			"slug" citext NOT NULL UNIQUE,
			"threads" integer DEFAULT 0,
			"title" text
		) WITH (autovacuum_enabled=false);

		CREATE UNLOGGED TABLE "threads"(
			"id" BIGSERIAL PRIMARY KEY,
			"author" citext REFERENCES users(nickname),
			"created" timestamp WITH TIME ZONE DEFAULT NOW(),
			"forum" citext NOT NULL REFERENCES forums(slug),
			"message" text NOT NULL,
			"slug" citext UNIQUE,
			"title" text,
			"votes" integer DEFAULT 0 NOT NULL
		) WITH (autovacuum_enabled=false);

		CREATE UNLOGGED TABLE "posts" (
			"id" BIGSERIAL,
			"parent" integer DEFAULT 0,
			"author" citext,
			"message" text NOT NULL,			
			"isedited" boolean DEFAULT 'false',
			"forum" citext REFERENCES forums(slug),
			"thread" bigint REFERENCES threads(id),
			"created" timestamp WITH TIME ZONE DEFAULT '1970-01-01T00:00:00.000Z',
			"path" bigint[]
		) WITH (autovacuum_enabled=false);

		CREATE UNLOGGED TABLE "votes" (
			"id" BIGSERIAL PRIMARY KEY,
			UNIQUE ("nickname", thread),
			"nickname" citext NOT NULL REFERENCES users(nickname),
			"thread" bigint REFERENCES threads(id),
			"voice" int2 
		) WITH (autovacuum_enabled=false);

		CREATE UNLOGGED TABLE "user_forum" (
			"nickname" citext COLLATE ucs_basic REFERENCES users(nickname),
			"forum" citext REFERENCES forums(slug),
			CONSTRAINT user_forum_const UNIQUE ("nickname", "forum")
		) WITH (autovacuum_enabled=false);

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

		-- posts:

		CREATE INDEX IF NOT EXISTS idxPostId 
			ON posts (id);

		CREATE INDEX IF NOT EXISTS idxPostPath 
			ON posts ((path[1]), path);

		CREATE INDEX IF NOT EXISTS idxPostThreadIdPath
			ON posts (thread, id, path);
		
		CREATE INDEX IF NOT EXISTS idxPostThreadPath
			ON posts (thread, path);

		-- sdfasdfs

		CREATE INDEX IF NOT EXISTS idxPostThreadParentId
			ON posts (path[1], id, parent)

		CREATE INDEX IF NOT EXISTS idxThreadAuthorForum
			ON threads (author, forum);

		-- ---------------------------------------------------------------------
		-- HELP FUNCTIONS:

		-- Threads:

		CREATE OR REPLACE FUNCTION threadInsert() RETURNS TRIGGER AS $thread_insert$
			BEGIN
				UPDATE forums 
				SET threads = threads + 1
				WHERE slug = NEW.forum;

				INSERT INTO user_forum (nickname, forum)
				VALUES (NEW.Author, NEW.Forum)
				ON CONFLICT DO NOTHING;
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

	// CREATE OR REPLACE FUNCTION postInsert() RETURNS TRIGGER AS $post_insert$
	// BEGIN
	// 	UPDATE forums
	// 	SET posts = posts + 1
	// 	WHERE slug = NEW.forum;
	// 	RETURN NULL;
	// END;
	// $post_insert$ LANGUAGE plpgsql;

	// CREATE TRIGGER postCreate
	// AFTER INSERT
	// ON posts
	// FOR EACH ROW EXECUTE PROCEDURE postInsert();

	// CREATE OR REPLACE FUNCTION postUpdate() RETURNS TRIGGER AS $post_update$
	// BEGIN
	// 	IF (NEW.message IS NOT NULL) AND (NEW.message <> OLD.message)
	// 		THEN
	// 			UPDATE posts
	// 			SET isedited = true
	// 			WHERE id = NEW.id;
	// 		END IF;
	// 	RETURN NEW;
	// END;
	// $post_update$ LANGUAGE plpgsql;

	// DROP TRIGGER IF EXISTS postUpdate ON posts;

	// CREATE TRIGGER postUpdate
	// BEFORE UPDATE
	// ON posts
	// FOR EACH ROW EXECUTE PROCEDURE postUpdate();

	// CREATE TRIGGER userForumInsert
	// AFTER INSERT
	// ON threads
	// FOR EACH ROW EXECUTE PROCEDURE userForumInsert();

	fmt.Println(sql)
	// _, err = DB.Exec(sql)

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
