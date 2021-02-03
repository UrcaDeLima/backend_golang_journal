package sqlstore

import (
	"database/sql"
	"log"

	"github.com/UrcaDeLima/backend_golang_journal/internal/app/model"
	"github.com/UrcaDeLima/backend_golang_journal/internal/app/store"
)

// NewsRepository ...
type NewsRepository struct {
	store *Store
}

// HeaderRepository ...
type HeaderRepository struct {
	store *Store
}

// ArticleRepository ...
type ArticleRepository struct {
	store *Store
}

// InnerDescriptionRepository ...
type InnerDescriptionRepository struct {
	store *Store
}

// PostRepository ...
type PostRepository struct {
	store *Store
}

// GetPostByID ...
func (r *PostRepository) GetPostByID(id int) (*model.PostModel, error) {
	postModel := &model.PostModel{}

	if err := r.store.db.QueryRow(
		"SELECT post.post_id, innerDescription.innerAdvertising, post.created_at, article.Title, article.backgroundImg, article.paragraphs, article.text, header.Title, header.image_id, header.date, header.views, header.shortDescription FROM post JOIN innerDescription ON post.innerAdvertising_id = innerDescription.innerDescription_id LEFT JOIN article ON post.post_id = article.post_id LEFT JOIN header ON post.post_id = header.post_id where post.post_id = $1",
		id,
	).Scan(
		&postModel.P.Post_id,
		&postModel.I.InnerAdvertising,
		&postModel.A.Title,
		&postModel.A.BackgroundImg,
		&postModel.A.Paragraphs,
		&postModel.P.Created_at,
		&postModel.A.Text,
		&postModel.H.Title,
		&postModel.H.Image_id,
		&postModel.H.Date,
		&postModel.H.Views,
		&postModel.H.ShortDescription,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}
	return postModel, nil
}

// GetAllPosts ...
func (r *PostRepository) GetAllPosts() ([]*model.PostModel, error) {
	allPostModel := []*model.PostModel{}

	rows, err := r.store.db.Query("SELECT post.post_id, innerDescription.innerAdvertising, post.created_at, article.Title, article.backgroundImg, article.paragraphs, article.text, header.Title, header.image_id, header.date, header.views, header.shortDescription FROM post JOIN innerDescription ON post.innerAdvertising_id = innerDescription.innerDescription_id LEFT JOIN article ON post.post_id = article.post_id LEFT JOIN header ON post.post_id = header.post_id ORDER BY created_at DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // спросить насчёт закрытия соединения

	for rows.Next() {
		postModel := &model.PostModel{}

		if err := rows.Scan(
			&postModel.P.Post_id,
			&postModel.P.Created_at,
			&postModel.I.InnerAdvertising,
			&postModel.A.Title,
			&postModel.A.BackgroundImg,
			&postModel.A.Paragraphs,
			&postModel.A.Text,
			&postModel.H.Title,
			&postModel.H.Image_id,
			&postModel.H.Date,
			&postModel.H.Views,
			&postModel.H.ShortDescription,
		); err != nil {
			log.Fatal(err)
		}
		allPostModel = append(allPostModel, postModel)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return allPostModel, err
}

// CreateNews ...
func (r *NewsRepository) CreateNews(news *model.News) error {
	return r.store.db.QueryRow(
		"INSERT INTO news (title, img, views, date) VALUES ($1, $2, $3, $4) RETURNING News_id",
		news.Title,
		news.Img,
		news.Views,
		news.Date,
	).Scan(&news.News_id)
}

// GetNewsByID ...
func (r *NewsRepository) GetNewsByID(id int) (*model.News, error) {
	u := &model.News{}

	if err := r.store.db.QueryRow(
		"SELECT News_id, Title, Img, Date, Views, Created_at FROM news WHERE News_id = $1",
		id,
	).Scan(
		&u.News_id,
		&u.Title,
		&u.Img,
		&u.Date,
		&u.Views,
		&u.Created_at,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// GetAllNews ...
func (r *NewsRepository) GetAllNews() ([]*model.News, error) {
	allNews := []*model.News{}

	rows, err := r.store.db.Query("SELECT News_id, Title, Img, Date, Views, Created_at FROM news ORDER BY created_at DESC LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		news := &model.News{}
		if err := rows.Scan(
			&news.News_id,
			&news.Title,
			&news.Img,
			&news.Date,
			&news.Views,
			&news.Created_at,
		); err != nil {
			log.Fatal(err)
		}
		allNews = append(allNews, news)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return allNews, err
}
