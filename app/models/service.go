package models

import (
	database "../db"
	slc "../sql"
)

type Service struct {
	Forum  int `json:"forum"`
	Post   int `json:"post"`
	Thread int `json:"thread"`
	User   int `json:"user"`
}

func (instance *Service) GetService() (err error) {

	err = database.DB.QueryRow(
		slc.GetService,
	).Scan(
		&instance.Forum,
		&instance.Post,
		&instance.User,
		&instance.Thread,
	)

	return
}

func (instance *Service) ClearDb() (err error) {
	_, err = database.DB.Exec(
		slc.ClearDb,
	)

	return
}
