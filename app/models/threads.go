package models

import (
	"strconv"
	"time"

	database "../db"
	slc "../sql"
)

type Thread struct {
	Id      int32      `json:"id"`
	Title   string     `json:"title"`
	Author  string     `json:"author"`
	Forum   string     `json:"forum"`
	Message string     `json:"message"`
	Votes   int32      `json:"votes"`
	Slug    *string    `json:"slug"`
	Created *time.Time `json:"created"`
}

type Threads []Thread

func (instance *Thread) CreateThread() (err error) {

	var u User

	// Проверяем наличие автора и форума
	err = database.DB.QueryRow(
		slc.GetUserByNickname,

		instance.Author,
	).Scan(
		&u.Nickname,
		&u.Fullname,
		&u.About,
		&u.Email,
	)

	if err != nil {
		return ThreadAuthorNotFoundError
	}

	// Проверям наличие форума
	err = database.DB.QueryRow(
		slc.CheckForumExistBySlug,

		instance.Forum,
	).Scan(
		&instance.Forum,
	)

	if err != nil {
		return ThreadForumNotFoundError
	}

	// Создаем тред
	threadErr := database.DB.QueryRow(
		slc.CreateThread,

		instance.Title,
		instance.Author,
		instance.Forum,
		instance.Message,
		instance.Slug,
		instance.Created,
	).Scan(
		&instance.Id,
	)

	// Пушим юзера в вспомогательную базу
	_, err = database.DB.Exec(
		slc.UserHelpTableInsert,

		u.Nickname,
		u.Fullname,
		u.About,
		u.Email,
		instance.Forum,
	)

	if err != nil {
		panic(err)
	}

	if threadErr == nil {
		return
	}

	err = database.DB.QueryRow(
		slc.GetThreadBySlug,

		instance.Slug,
	).Scan(
		&instance.Id,
		&instance.Title,
		&instance.Author,
		&instance.Forum,
		&instance.Message,
		&instance.Votes,
		&instance.Slug,
		&instance.Created,
	)

	if err == nil {
		return ThreadAlreadyExistError
	}

	return
}

func (instance *Thread) GetThread(soi string) (err error) {
	// Парсим slug/url
	val, parseErr := strconv.Atoi(soi)

	if parseErr != nil {
		err = database.DB.QueryRow(
			slc.GetThreadBySlug,

			soi,
		).Scan(
			&instance.Id,
			&instance.Title,
			&instance.Author,
			&instance.Forum,
			&instance.Message,
			&instance.Votes,
			&instance.Slug,
			&instance.Created,
		)
	} else {
		err = database.DB.QueryRow(
			slc.GetThreadById,

			val,
		).Scan(
			&instance.Id,
			&instance.Title,
			&instance.Author,
			&instance.Forum,
			&instance.Message,
			&instance.Votes,
			&instance.Slug,
			&instance.Created,
		)
	}

	if err != nil {
		return ThreadNotFoundError
	}

	return
}

func (instance *Thread) ThreadUpdate(soi string) (err error) {

	// Парсим slug/url
	val, parseErr := strconv.Atoi(soi)

	if parseErr != nil {
		err = database.DB.QueryRow(
			slc.GetThreadSlugAndIdBySlug,

			soi,
		).Scan(
			&instance.Forum,
			&instance.Id,
		)
	} else {
		err = database.DB.QueryRow(
			slc.GetThreadSlugAndIdByID,

			val,
		).Scan(
			&instance.Forum,
			&instance.Id,
		)
	}

	if err != nil {
		return ThreadNotFoundError
	}

	err = database.DB.QueryRow(
		slc.ThreadUpdate,

		instance.Id,
		instance.Message,
		instance.Title,
	).Scan(
		&instance.Forum,
		&instance.Author,
		&instance.Slug,
		&instance.Created,
		&instance.Message,
		&instance.Title,
		&instance.Votes,
	)

	return
}
