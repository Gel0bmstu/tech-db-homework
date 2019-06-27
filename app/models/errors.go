package models

import "errors"

// Users errors
var (
	UserConflictError = errors.New("User already exixst")
	UserNotFoundError = errors.New("User not found")
)

// Forunms errors
var (
	ForumUserNotFoundError = errors.New("Can't find user")
	ForumConflictError     = errors.New("Forum already exist")
	ForumNotFoundError     = errors.New("Forum not found")
)

// Threads error
var (
	ThreadAuthorNotFoundError = errors.New("Can't find user")
	ThreadForumNotFoundError  = errors.New("Can't find forum")
	ThreadAlreadyExistError   = errors.New("Thread already exist")
	ThreadNotFoundError       = errors.New("Thread not found")
)

// Posts error
var (
	ParentPostNotFoundInThread = errors.New("Parent post not found in thread")
	PostNotFoundError          = errors.New("Post not found")
)
