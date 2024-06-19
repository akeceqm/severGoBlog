package services

import (
	"github.com/jmoiron/sqlx"
	"post/internal/database/models"
	"time"
)

func CreatePost(title, text string, authorId string, db *sqlx.DB) (models.Post, error) {
	post := models.Post{}

	err := db.Get("SELECT * FROM public.post WHERE id = $1", authorId)

	if err != nil {
		return models.Post{}, err
	}

	post.Id = GenerateId()
	post.Title = title
	post.Text = text
	post.DateCreated = time.Now()
	post.AuthorId = authorId

	_, err = db.Exec("INSERT INTO public.post (id, title, text, date, author_id) VALUES ($1, $2, $3, $4, $5)", post.Id, post.Title, post.Text, post.DateCreated, post.AuthorId)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func GetPosts(db *sqlx.DB) ([]models.Post, error) {
	post := []models.Post{}

	err := db.Select(&post, "SELECT * FROM public.post")
	if err != nil {
		return post, err
	}
	return post, nil
}

func GetPostById(id string, db *sqlx.DB) (models.Post, error) {
	post := models.Post{}

	err := db.Get(&post, "SELECT * FROM public.post WHERE id = $1", id)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func GetPostFull(db *sqlx.DB) ([]models.FullPost, error) {
	post := []models.FullPost{}

	err := db.Select(&post, "SELECT public.post.*, public.user.nick_name FROM public.post JOIN public.user ON public.post.author_id = public.user.id")
	if err != nil {
		return post, err
	}
	return post, nil
}

func GetPostsByAuthorId(idAuthor string, db *sqlx.DB) ([]models.Post, error) {
	var post []models.Post

	err := db.Select(&post, "SELECT * FROM public.post WHERE author_id = $1", idAuthor)
	if err != nil {
		return []models.Post{}, err
	}
	return post, nil
}

func GetPostByBetweenDate(startDate, endDate time.Time, db *sqlx.DB) ([]models.Post, error) {
	var post []models.Post

	err := db.Select(&post, "SELECT * FROM public.post WHERE date_created BETWEEN $1 AND $2", startDate, endDate)
	if err != nil {
		return []models.Post{}, err
	}

	return post, nil
}

func UpdatePost(id string, title, text string, db *sqlx.DB) error {
	var postDB models.Post

	err := db.Get(&postDB, "SELECT * FROM public.post WHERE id = $1", id)
	if err != nil {
		return err
	}
	if title == "" {
		title = postDB.Title
	}
	if text == "" {
		text = postDB.Text
	}

	_, err = db.Exec("UPDATE public.post SET title = $1, text = $2 WHERE id = $3", title, text, id)
	if err != nil {
		return err
	}

	return nil
}

func DeletePost(id string, db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM public.post WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func GetPostsOrder(number int, db *sqlx.DB) ([]models.FullPost, error) {
	var posts []models.FullPost

	err := db.Select(&posts, "SELECT public.post.*, public.user.nick_name, COALESCE(comment_counts.comment_count, 0) AS comment_count FROM public.post JOIN public.user ON public.post.author_id = public.user.id LEFT JOIN (SELECT post_id, COUNT(*) AS comment_count FROM public.comment GROUP BY post_id) AS comment_counts ON public.post.id = comment_counts.post_id ORDER BY public.post.date_created DESC LIMIT 10 OFFSET $1", (number-1)*10)
	if err != nil {
		return []models.FullPost{}, err
	}

	var FullPosts []models.FullPost
	for i := 0; i < len(posts); i++ {
		FullPosts = append(FullPosts, models.FullPost{
			Id:                posts[i].Id,
			Title:             posts[i].Title,
			Text:              posts[i].Text,
			AuthorId:          posts[i].AuthorId,
			DateCreatedFormat: posts[i].DateCreated.Format("2006-01-02 15:04:05"),
			AuthorName:        posts[i].AuthorName,
			CommentsCount:     posts[i].CommentsCount,
		})
	}

	return FullPosts, nil
}
