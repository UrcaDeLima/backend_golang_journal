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
func (r *PostRepository) GetPostByID(id int) (*model.Post, *model.Header, *model.Article, *model.InnerDescription, error) {
	p := &model.Post{}
	h := &model.Header{}
	a := &model.Article{}
	i := &model.InnerDescription{}

	if err := r.store.db.QueryRow(
		"SELECT post.post_id, innerDescription.innerAdvertising, post.created_at, article.Title, article.backgroundImg, article.paragraphs, article.text, header.Title, header.image_id, header.date, header.views, header.shortDescription FROM post JOIN innerDescription ON post.innerAdvertising_id = innerDescription.innerDescription_id LEFT JOIN article ON post.post_id = article.post_id LEFT JOIN header ON post.post_id = header.post_id where post.post_id = $1",
		id,
	).Scan(
		&p.Post_id,
		&i.InnerAdvertising,
		&a.Title,
		&a.BackgroundImg,
		&a.Paragraphs,
		&p.Created_at,
		&a.Text,
		&h.Title,
		&h.Image_id,
		&h.Date,
		&h.Views,
		&h.ShortDescription,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil, store.ErrRecordNotFound
		}

		return nil, nil, nil, nil, err
	}

	return p, h, a, i, nil
}

// GetAllPosts ...
func (r *PostRepository) GetAllPosts() ([]*model.Post, []*model.Header, []*model.Article, []*model.InnerDescription, error) {
	allPosts := []*model.Post{}
	allHeaders := []*model.Header{}
	allArticles := []*model.Article{}
	allInnerDescriptions := []*model.InnerDescription{}

	rows, err := r.store.db.Query("SELECT post.post_id, innerDescription.innerAdvertising, post.created_at, article.Title, article.backgroundImg, article.paragraphs, article.text, header.Title, header.image_id, header.date, header.views, header.shortDescription FROM post JOIN innerDescription ON post.innerAdvertising_id = innerDescription.innerDescription_id LEFT JOIN article ON post.post_id = article.post_id LEFT JOIN header ON post.post_id = header.post_id ORDER BY created_at DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // спросить насчёт закрытия соединения

	for rows.Next() {
		p := &model.Post{}
		h := &model.Header{}
		a := &model.Article{}
		i := &model.InnerDescription{}
		if err := rows.Scan(
			&p.Post_id,
			&p.Created_at,
			&i.InnerAdvertising,
			&a.Title,
			&a.BackgroundImg,
			&a.Paragraphs,
			&a.Text,
			&h.Title,
			&h.Image_id,
			&h.Date,
			&h.Views,
			&h.ShortDescription,
		); err != nil {
			log.Fatal(err)
		}
		allPosts = append(allPosts, p)
		allHeaders = append(allHeaders, h)
		allArticles = append(allArticles, a)
		allInnerDescriptions = append(allInnerDescriptions, i)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return allPosts, allHeaders, allArticles, allInnerDescriptions, err
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
