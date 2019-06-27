package router

import (
	"../handlers"
	"github.com/gin-gonic/gin"
)

func SetRoutes(instance *gin.Engine) {

	r := instance.Group("/api")
	{
		// GET
		r.GET("/forum/:slug/details", handlers.ForumSlugDetails)
		r.GET("/forum/:slug/users", handlers.ForumSlugUsers)
		r.GET("/forum/:slug/threads", handlers.ForumSlugThreads)

		r.GET("/post/:id/details", handlers.PostIdDetailsGet)

		r.GET("/service/status", handlers.ServiceStatus)

		r.GET("/thread/:slugOrID/details", handlers.ThreadsSlugOrIdDetailsGet)
		r.GET("/thread/:slugOrID/posts", handlers.ThreadsSlugOrIdPosts)

		r.GET("/user/:nickname/profile", handlers.UserNicknameProfileGet)

		// POST
		r.POST("/forum/:slug", handlers.ForumCreate)
		r.POST("/forum/:slug/create", handlers.ForumSlugCreate)

		r.POST("/post/:id/details", handlers.PostIdDetailsPost)

		r.POST("/service/clear", handlers.ServiceClear)

		r.POST("/thread/:slugOrID/create", handlers.ThreadsSlugOrIdCrete)
		r.POST("/thread/:slugOrID/details", handlers.ThreadsSlugOrIdDetailsPost)
		r.POST("/thread/:slugOrID/vote", handlers.ThreadsSlugOrIdVote)

		r.POST("/user/:nickname/create", handlers.UserNicknameCreate)
		r.POST("/user/:nickname/profile", handlers.UserNicknameProfilePost)
	}

}

func Initialize() (r *gin.Engine) {

	r = gin.Default()

	SetRoutes(r)

	return
}
