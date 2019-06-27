package models

import (
	"time"

	database "../db"
	slc "../sql"
	"github.com/jackc/pgx"
)

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts"`
	Threads int32  `json:"threads"`
}

type UrlVars struct {
	Limit string
	Since string
	Sort  string
	Desc  string
}

func (instance *Forum) CreateForum() (err error) {

	// Проверяем наличие пользователя в базе
	err = database.DB.QueryRow(
		slc.CheckUserNicknameByNickname,

		instance.User,
	).Scan(
		&instance.User,
	)

	if err != nil {
		return ForumUserNotFoundError
	}

	// Проверяем наличие форума в базе
	err = database.DB.QueryRow(
		slc.GetForumBySlug,

		instance.Slug,
	).Scan(
		&instance.Title,
		&instance.User,
		&instance.Slug,
		&instance.Posts,
		&instance.Threads,
	)

	if err == nil {
		return ForumConflictError
	}

	// Создаем форум
	_, err = database.DB.Exec(
		slc.CreateForum,

		instance.Title,
		instance.User,
		instance.Slug,
	)

	return
}

func (instance *Forum) GetForumDetails() (err error) {
	err = database.DB.QueryRow(
		slc.GetForumBySlug,

		instance.Slug,
	).Scan(
		&instance.Title,
		&instance.User,
		&instance.Slug,
		&instance.Posts,
		&instance.Threads,
	)

	if err != nil {
		return ForumNotFoundError
	}

	return
}

func (instance *Forum) GetThreads(uv UrlVars) (threads Threads, err error) {
	var (
		rows         *pgx.Rows
		RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	)

	err = database.DB.QueryRow(
		slc.CheckForumExistBySlug,

		instance.Slug,
	).Scan(
		&instance.Slug,
	)

	if err != nil {
		return nil, ForumNotFoundError
	}

	if uv.Since != "" {

		since, err := time.Parse(RFC3339Milli, string(uv.Since))

		if err != nil {
			return nil, err
		}

		if uv.Desc == "DESC" {
			rows, err = database.DB.Query(
				slc.GetThreadsByForumSinceDESC,

				instance.Slug,
				since,
				uv.Limit,
			)
		} else {
			rows, err = database.DB.Query(
				slc.GetThreadsByForumSince,

				instance.Slug,
				since,
				uv.Limit,
			)
		}
	} else {
		if uv.Desc == "DESC" {
			rows, err = database.DB.Query(
				slc.GetThreadsByForumDESC,

				instance.Slug,
				uv.Limit,
			)
		} else {
			rows, err = database.DB.Query(
				slc.GetThreadsByForum,

				instance.Slug,
				uv.Limit,
			)
		}
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t Thread
		if err := rows.Scan(
			&t.Id,
			&t.Title,
			&t.Forum,
			&t.Author,
			&t.Message,
			&t.Slug,
			&t.Created,
			&t.Votes,
		); err != nil {
			return nil, err
		}

		threads = append(threads, t)
	}

	return
}
