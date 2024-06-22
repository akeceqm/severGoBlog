package handler

import (
	"log"
	"net/http"
	"post/cmd"
	"post/internal/database/models"
	"post/internal/handler/handlerUser"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"post/internal/handler/handlerComment"
	"post/internal/handler/handlerPost"
	"post/internal/services"
)

func InitRoutesHTML(server *gin.Engine, db *sqlx.DB) {
	authMiddleware := AuthMiddleware(db)

	server.GET("/authorization", func(c *gin.Context) {
		c.HTML(200, "authorization.html", gin.H{})
	})

	server.GET("/registration", func(c *gin.Context) {
		c.HTML(200, "registration.html", gin.H{})
	})
	// Применяем middleware авторизации
	server.Use(authMiddleware)

	cmd.Server.GET("/", func(c *gin.Context) {
		handlerIndex(db, c)

	})

	server.GET("/profileUser", func(c *gin.Context) {
		c.HTML(200, "profileUser.html", gin.H{})
	})

	server.GET("/profileUser/:userId", func(c *gin.Context) {
		c.HTML(200, "profileUser.html", gin.H{})
	})
	server.GET("/changeProfile", func(c *gin.Context) {
		c.HTML(200, "changeProfile.html", gin.H{})
	})
	server.GET("/changeProfile/:userId", func(c *gin.Context) {
		user, err := handlerUser.PUTHandleUser(c, db)
		if err != nil {
			c.HTML(400, "400.html", gin.H{"Error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	})
	server.GET("/h/post/:idPost/comments", func(c *gin.Context) {
		handlerComment.GETHandlePostCommentsHTML(c, db)
	})
	server.GET("/h/:countPage", func(c *gin.Context) {
		handlerPost.GETHandlePostsHTML(c, db)
	})

	server.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})
}

func handlerIndex(db *sqlx.DB, c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists || userID == nil {
		c.HTML(400, "400.html", gin.H{"Error": "Пользователь не авторизован или сессия истекла"})
		handlerIndexNoAuthorization(c, db)
		return
	}

	// Проверка авторизации
	isAuthorized, err := services.IsUserAuthorized(db, userID.(string))
	if err != nil {
		log.Println("Ошибка проверки авторизации:", err)
		handlerIndexNoAuthorization(c, db)
		return
	}

	if isAuthorized {
		handlerIndexAuthorization(c)
	} else {
		handlerIndexNoAuthorization(c, db)
	}
}

func handlerIndexNoAuthorization(c *gin.Context, db *sqlx.DB) {
	log.Println("Rendering PageMainNoAuthorization.html")
	post, err := services.GetPostFull(db)
	if err != nil {
		c.HTML(400, "400.html", gin.H{"Error": err.Error()})
		return
	}

	if len(post) == 0 {
		log.Println("No posts found")
		c.HTML(200, "PageMainNoAuthorization.html", gin.H{"posts": []models.FullPost{}})
		return
	}

	var fullPosts []models.FullPost
	for i := 0; i < 10 && i < len(post); i++ {
		comments, err := services.GetCommentsByPostId(post[i].Id, db)
		if err != nil {
			c.HTML(400, "400.html", gin.H{"Error": err.Error()})
			return
		}

		fullPosts = append(fullPosts, models.FullPost{
			Id:                post[i].Id,
			Title:             post[i].Title,
			Text:              post[i].Text,
			AuthorId:          post[i].AuthorId,
			DateCreatedFormat: post[i].DateCreated.Format("2006-01-02 15:04:05"),
			AuthorName:        post[i].AuthorName,
			Comments:          []models.FullComment{},
			CommentsCount:     len(comments),
		})
	}
	c.HTML(200, "PageMainNoAuthorization.html", gin.H{"posts": fullPosts})
}

func handlerIndexAuthorization(c *gin.Context) {
	log.Println("Rendering PagePostComments.html")
	c.HTML(200, "PageMainYesAuthorization.html", nil)
}

func AuthMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			handlerIndexNoAuthorization(c, db)
			c.Abort()
			return
		}

		session, err := services.GetSessionByID(db, sessionID)
		if err != nil || session.UserID == "" {
			handlerIndexNoAuthorization(c, db)
			c.Abort()
			return
		}

		c.Set("userID", session.UserID) // Установка userID в контекст Gin для авторизованных пользователей
		c.Next()
	}
}
