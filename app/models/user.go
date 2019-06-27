package models

import (
	database "../db"
	slc "../sql"
	"github.com/jackc/pgx"
)

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type Users []User

func (instance *User) CreateUser() (users Users, err error) {

	// Создаем пользователя
	_, err = database.DB.Exec(
		slc.CreateUser,

		instance.Nickname,
		instance.Fullname,
		instance.About,
		instance.Email,
	)

	// Если есть конфликтующие пользователи, помещаем их в массив users
	if err != nil {
		conflictUsers, err := database.DB.Query(
			slc.GetUserByNicknameOrEmail,

			instance.Nickname,
			instance.Email,
		)

		if err != nil {
			return nil, err
		}

		for conflictUsers.Next() {
			var u User
			err := conflictUsers.Scan(
				&u.Nickname,
				&u.Fullname,
				&u.About,
				&u.Email,
			)

			if err != nil {
				return nil, err
			}

			users = append(users, u)
		}

		return users, UserConflictError
	}

	return
}

func (instance *User) EditUserInfo() (err error) {

	// Проверяем наличие пользователя по никнейму
	err = database.DB.QueryRow(
		slc.CheckUserNicknameByNickname,

		instance.Nickname,
	).Scan(
		&instance.Nickname,
	)

	if err != nil {
		return UserNotFoundError
	}

	// Проверяем наличие конфликтующих мейлов
	err = database.DB.QueryRow(
		slc.CheckUserEmailByEmail,

		instance.Email,
	).Scan(
		&instance.Email,
	)

	if err == nil {
		return UserConflictError
	}

	// Обновляем информацию о пользователе
	err = database.DB.QueryRow(
		slc.UpdateUser,

		instance.Nickname,
		instance.Fullname,
		instance.About,
		instance.Email,
	).Scan(
		&instance.Fullname,
		&instance.About,
		&instance.Email,
	)

	return
}

func (instance *User) GetUserInfo() (err error) {
	// Достаем информацию о пользователе по никнейму или мейлу
	err = database.DB.QueryRow(
		slc.GetUserByNicknameOrEmail,

		instance.Nickname,
		instance.Email,
	).Scan(
		&instance.Nickname,
		&instance.Fullname,
		&instance.About,
		&instance.Email,
	)

	if err != nil {
		return UserNotFoundError
	}

	return
}

func (instance *Users) GetUsers(uv UrlVars, s string) (err error) {

	var (
		rows *pgx.Rows
	)

	err = database.DB.QueryRow(
		slc.CheckForumExistBySlug,

		s,
	).Scan(
		&s,
	)

	if err != nil {
		return ForumNotFoundError
	}

	// пофиксить какие-то странные костыли с лимитами
	if uv.Since != "" {
		if uv.Desc == "DESC" {
			if uv.Limit != "" {
				rows, err = database.DB.Query(
					slc.GetForumUsersSinceDesc+"LIMIT $3::TEXT::INTEGER",

					s,
					uv.Since,
					uv.Limit,
				)
			} else {
				rows, err = database.DB.Query(
					slc.GetForumUsersSinceDesc+";",

					s,
					uv.Since,
				)
			}
		} else {
			if uv.Limit != "" {
				rows, err = database.DB.Query(
					slc.GetForumUsersSince+"LIMIT $3::TEXT::INTEGER",

					s,
					uv.Since,
					uv.Limit,
				)
			} else {
				rows, err = database.DB.Query(
					slc.GetForumUsersSince+";",

					s,
					uv.Since,
				)
			}
		}
	} else {
		if uv.Desc == "DESC" {
			if uv.Limit != "" {
				rows, err = database.DB.Query(
					slc.GetForumUsersDesc+"LIMIT $2::TEXT::INTEGER",

					s,
					uv.Limit,
				)
			} else {
				rows, err = database.DB.Query(
					slc.GetForumUsersDesc+";",

					s,
				)
			}
		} else {
			if uv.Limit != "" {
				rows, err = database.DB.Query(
					slc.GetForumUsers+"LIMIT $2::TEXT::INTEGER",

					s,
					uv.Limit,
				)
			} else {
				rows, err = database.DB.Query(
					slc.GetForumUsers+";",

					s,
				)
			}
		}
	}

	for rows.Next() {
		var u User

		err = rows.Scan(
			&u.Nickname,
			&u.Fullname,
			&u.About,
			&u.Email,
		)

		if err != nil {
			return
		}

		*instance = append(*instance, u)
	}

	return
}
