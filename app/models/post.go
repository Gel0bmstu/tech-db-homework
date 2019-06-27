package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx"

	database "../db"
	slc "../sql"

	"strconv"
)

type Post struct {
	Id       int64      `json:"id"`
	Parent   int64      `json:"parent"`
	Author   string     `json:"author"`
	Message  string     `json:"message"`
	IsEdited bool       `json:"isEdited"`
	Forum    *string    `json:"forum"`
	Thread   *int32     `json:"thread"`
	Created  *time.Time `json:"created"`
	Path     []string
}

type Posts []Post

func (instance *Posts) CreatePost(soi string) (err error) {

	var (
		threadID  int32
		forumSlug string

		queryPosts     string
		queryUserForum string

		postsInsertValues     []interface{}
		userForumInsertValues []interface{}
		postCount             int
	)

	val, parseErr := strconv.Atoi(soi)

	// Парсим слаг/ид
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

	if err != nil {
		return ThreadNotFoundError
	}

	if len(*instance) == 0 {
		return
	}

	// Задаем квери для массового инсерта в таблицу
	queryPosts = `INSERT INTO posts
				  (author, message, parent, thread, forum)
				  VALUES`

	queryUserForum = `INSERT INTO user_forum
					 (nickname, forum)
					 VALUES`

	for i, _ := range *instance {

		(*instance)[i].Forum = &forumSlug
		(*instance)[i].Thread = &threadID

		if (*instance)[i].Author != "" {
			err = database.DB.QueryRow(
				slc.CheckUserNicknameByNickname,

				(*instance)[i].Author,
			).Scan(
				&(*instance)[i].Author,
			)

			if err != nil {
				return UserNotFoundError
			}
		}

		if (*instance)[i].Parent != 0 {
			err = database.DB.QueryRow(
				slc.CheckExistingPostByThreadId,

				(*instance)[i].Thread,
				(*instance)[i].Parent,
			).Scan(
				&(*instance)[i].Parent,
			)

			if err != nil {
				return ParentPostNotFoundInThread
			}
		}

		postsInsertValues = append(postsInsertValues, (*instance)[i].Author, (*instance)[i].Message, (*instance)[i].Parent, (*instance)[i].Thread, (*instance)[i].Forum)
		userForumInsertValues = append(userForumInsertValues, (*instance)[i].Author, (*instance)[i].Forum)

		qPosts := fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d)",

			i*5+1,
			i*5+2,
			i*5+3,
			i*5+4,
			i*5+5,
		)

		qUserForum := fmt.Sprintf(
			"($%d, $%d)",

			i*2+1,
			i*2+2,
		)

		queryPosts += qPosts
		queryUserForum += qUserForum

		if i != len(*instance)-1 {
			queryPosts += ", "
			queryUserForum += ", "
		}
	}

	queryPosts += ` RETURNING created, id;`
	queryUserForum += ` ON CONFLICT DO NOTHING;`

	// Заполняем вспомогательную таблицу юзеров
	transactionUserFroum, _ := database.DB.Begin()

	_, err = transactionUserFroum.Exec(
		queryUserForum,

		userForumInsertValues...,
	)

	transactionUserFroum.Commit()

	// Заполняем посты недостающими значениями
	transactionPosts, _ := database.DB.Begin()

	rows, err := transactionPosts.Query(
		queryPosts,

		postsInsertValues...,
	)

	if err != nil {
		fmt.Println("POLUCHAY PIDOR!")
		transactionPosts.Rollback()
		return
	}

	postCount = 0
	for rows.Next() {
		rows.Scan(
			&(*instance)[postCount].Created,
			&(*instance)[postCount].Id,
		)
		postCount++
	}

	if err != nil {
		fmt.Println("OTSOSI!")
		transactionPosts.Rollback()
		return
	}

	transactionPosts.Commit()

	transactionForum, _ := database.DB.Begin()
	_, err = transactionForum.Exec(slc.ForumPostsCoutnUpdate, forumSlug, postCount)

	if err != nil {
		fmt.Println("NA, POLUCHAY!")
		transactionForum.Rollback()
		return
	}

	transactionForum.Commit()
	return
}

func (instance *Posts) GetPosts(soi string, uv UrlVars) (err error) {
	var (
		threadID  int32
		forumSlug string
		rows      *pgx.Rows
	)

	val, parseErr := strconv.Atoi(soi)

	// Парсим слаг/ид
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

	if err != nil {
		return ThreadNotFoundError
	}

	switch uv.Sort {
	case "flat", "":
		rows, err = GetPostsByFlatSort(uv, threadID)
	case "tree":
		rows, err = GetPostsByTreeSort(uv, threadID)
	case "parent_tree":
		rows, err = GetPostsByParentTreeSort(uv, threadID)
	}

	if err != nil {
		return err
	}

	for rows.Next() {

		var p Post

		p.Forum = &forumSlug

		err = rows.Scan(
			&p.Id,
			&p.Author,
			&p.Message,
			&p.Thread,
			&p.Created,
			&p.Parent,
		)

		if err != nil {
			return err
		}

		*instance = append(*instance, p)

	}

	return
}

func GetPostsByFlatSort(uv UrlVars, threadID int32) (rows *pgx.Rows, err error) {

	if uv.Since != "" {
		if uv.Desc == "DESC" {
			rows, err = database.DB.Query(
				slc.GetPostsByIdFlatSinceDesc,

				threadID,
				uv.Since,
				uv.Limit,
			)
		} else {
			rows, err = database.DB.Query(
				slc.GetPostsByIdFlatSince,

				threadID,
				uv.Since,
				uv.Limit,
			)
		}
	} else {
		if uv.Desc == "DESC" {
			rows, err = database.DB.Query(
				slc.GetPostsByIdFlatDesc,

				threadID,
				uv.Limit,
			)
		} else {
			rows, err = database.DB.Query(
				slc.GetPostsByIdFlat,

				threadID,
				uv.Limit,
			)
		}
	}

	return
}

func GetPostsByTreeSort(uv UrlVars, threadID int32) (rows *pgx.Rows, err error) {

	if uv.Since != "" {
		if uv.Desc == "DESC" {
			rows, err = database.DB.Query(
				slc.GetPostsByIdTreeSinceDesc,

				threadID,
				uv.Since,
				uv.Limit,
			)
		} else {
			rows, err = database.DB.Query(
				slc.GetPostsByIdTreeSince,

				threadID,
				uv.Since,
				uv.Limit,
			)
		}
	} else {
		if uv.Desc == "DESC" {
			rows, err = database.DB.Query(
				slc.GetPostsByIdTreeDesc,

				threadID,
				uv.Limit,
			)
		} else {
			rows, err = database.DB.Query(
				slc.GetPostsByIdTree,

				threadID,
				uv.Limit,
			)
		}
	}

	return
}

func GetPostsByParentTreeSort(uv UrlVars, threadID int32) (rows *pgx.Rows, err error) {

	if uv.Since != "" {
		if uv.Desc == "DESC" {
			rows, err = database.DB.Query(
				slc.GetPostsByIdParentTreeSinceDesc,

				threadID,
				uv.Since,
				uv.Limit,
			)
		} else {
			rows, err = database.DB.Query(
				slc.GetPostsByIdParentTreeSince,

				threadID,
				uv.Since,
				uv.Limit,
			)
		}
	} else {
		if uv.Desc == "DESC" {
			rows, err = database.DB.Query(
				slc.GetPostsByIdParentTreeDesc,

				threadID,
				uv.Limit,
			)
		} else {
			rows, err = database.DB.Query(
				slc.GetPostsByIdParentTree,

				threadID,
				uv.Limit,
			)
		}
	}

	return
}

func (instance *Post) GetPostDetails(id string, r string) (data map[string]interface{}, err error) {

	data = make(map[string]interface{})

	err = database.DB.QueryRow(
		slc.GetPostById,

		id,
	).Scan(
		&instance.Id,
		&instance.Author,
		&instance.Message,
		&instance.Thread,
		&instance.Created,
		&instance.Parent,
		&instance.Forum,
		&instance.IsEdited,
	)

	if err != nil {
		return nil, PostNotFoundError
	}

	data["post"] = instance

	if r == "" {
		return
	}

	rVals := strings.Split(r, ",")

	for _, val := range rVals {
		switch val {
		case "user":
			var u User
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
				return nil, UserNotFoundError
			}
			data["author"] = u
		case "forum":
			var f Forum
			err = database.DB.QueryRow(
				slc.GetForumBySlug,

				instance.Forum,
			).Scan(
				&f.Title,
				&f.User,
				&f.Slug,
				&f.Posts,
				&f.Threads,
			)
			if err != nil {
				return nil, ForumNotFoundError
			}
			data["forum"] = f
		case "thread":
			var t Thread
			err = database.DB.QueryRow(
				slc.GetThreadById,

				instance.Thread,
			).Scan(
				&t.Id,
				&t.Title,
				&t.Author,
				&t.Forum,
				&t.Message,
				&t.Votes,
				&t.Slug,
				&t.Created,
			)
			if err != nil {
				return nil, ThreadNotFoundError
			}

			data["thread"] = t
		}
	}

	return
}

func (instance *Post) PostUpdate(id string) (err error) {

	var oldMessage string

	err = database.DB.QueryRow(
		slc.CheckPostByIdAndGetMessege,

		id,
	).Scan(
		&instance.Id,
		&oldMessage,
	)

	if err != nil {
		return PostNotFoundError
	}

	if instance.Message == oldMessage || instance.Message == "" {
		err = database.DB.QueryRow(
			slc.GetPostById,

			id,
		).Scan(
			&instance.Id,
			&instance.Author,
			&instance.Message,
			&instance.Thread,
			&instance.Created,
			&instance.Parent,
			&instance.Forum,
			&instance.IsEdited,
		)
	} else {
		err = database.DB.QueryRow(
			slc.UpdatePostById,

			id,
			instance.Message,
		).Scan(
			&instance.Id,
			&instance.Author,
			&instance.Message,
			&instance.Thread,
			&instance.Created,
			&instance.Parent,
			&instance.Forum,
			&instance.IsEdited,
		)
	}

	return
}
