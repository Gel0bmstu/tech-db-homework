package models

import (
	"strconv"

	database "../db"
	slc "../sql"
)

type Vote struct {
	Id       int
	Thread   int32
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}

func (instanse *Vote) ChangeVote(soi string) (*Thread, error) {

	var (
		threadID  int32
		forumSlug string
		thread    Thread
	)

	// Пороверяем наличие юзера
	err := database.DB.QueryRow(
		slc.CheckUserNicknameByNickname,

		instanse.Nickname,
	).Scan(
		&instanse.Nickname,
	)

	if err != nil {
		return nil, UserNotFoundError
	}

	// Парсим slug/url
	val, parseErr := strconv.Atoi(soi)

	if parseErr != nil {
		err = database.DB.QueryRow(
			slc.GetThreadSlugAndIdBySlug,

			soi,
		).Scan(
			&forumSlug,
			&threadID,
		)
	} else {
		err = database.DB.QueryRow(
			slc.GetThreadSlugAndIdByID,

			val,
		).Scan(
			&forumSlug,
			&threadID,
		)
	}

	instanse.Thread = threadID

	if err != nil {
		return nil, ThreadNotFoundError
	}

	// Создаем голос
	_, err = database.DB.Exec(
		slc.CreateVote,

		instanse.Thread,
		instanse.Nickname,
		instanse.Voice,
	)

	if err != nil {
		return nil, err
	}

	// Получаем информациб о треде
	err = database.DB.QueryRow(
		slc.GetThreadById,

		threadID,
	).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)

	if err != nil {
		return nil, err
	}

	return &thread, err
}
