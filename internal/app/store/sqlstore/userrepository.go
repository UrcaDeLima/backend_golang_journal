package sqlstore

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

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

// ImageRepository ...
type ImageRepository struct {
	store *Store
}

// PostRepository ...
type PostRepository struct {
	store *Store
}

// UpdatePicture ...
func (r *ImageRepository) UpdatePicture(id int, m *multipart.Reader) error {
	image := model.Image{}

	for {
		part, err := m.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			continue
		}

		dst, err := os.Create("./image/" + part.FileName())
		if err != nil {
			fmt.Println(err)
		}

		if part.FormName() == "desktop" {
			image.Desktop = "image/" + part.FileName()
		} else if part.FormName() == "mobile" {
			image.Mobile = "image/" + part.FileName()
		}

		// image, err := jpeg.DecodeConfig(part)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "%s: %v\n", "backend_golang_journal/image/"+part.FileName(), err)
		// }
		// log.Println(image)

		io.Copy(dst, part)
	}
	log.Println(image)
	return r.store.db.QueryRow(
		"UPDATE image SET desktop = $1, mobile = $2 WHERE image_id = $3 RETURNING image_id",
		image.Desktop,
		image.Mobile,
		id,
	).Scan(&image.Image_id)
}

// SetPicture ...
func (r *ImageRepository) SetPicture(m *multipart.Reader) error {
	image := *&model.Image{}

	for {
		part, err := m.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			continue
		}

		// creating pictures
		dst, err := os.Create("./image/" + part.FileName())
		if err != nil {
			fmt.Println(err)
		}

		// part query for sql
		if part.FormName() == "desktop" {
			image.Desktop = "image/" + part.FileName()
		} else if part.FormName() == "mobile" {
			image.Mobile = "image/" + part.FileName()
		}

		io.Copy(dst, part)

		///////////////////// to get a size
		// image, err := jpeg.DecodeConfig(part)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "%s: %v\n", "image/"+part.FileName(), err)
		// }
		// log.Println(image)

		/////////////////// to change a size
		// // load original image
		// img, err := imaging.Open("image/" + part.FileName())
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }

		// dstimg := imaging.Resize(img, 1920, 0, imaging.Box)

		// // save resized image
		// err = imaging.Save(dstimg, "image/"+part.FileName())

		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }

		// // everything ok
		// fmt.Println("Image resized")
	}
	return r.store.db.QueryRow(
		"INSERT INTO image (desktop, mobile) VALUES ($1, $2) RETURNING image_id",
		image.Desktop,
		image.Mobile,
	).Scan(&image.Image_id)
}

// GetPostByID ...
func (r *PostRepository) GetPostByID(id int) (*model.PostModel, error) {
	postModel := &model.PostModel{}

	if err := r.store.db.QueryRow(
		"SELECT post.post_id, innerDescription.innerAdvertising, post.created_at, article.Title, article.backgroundImg, article.paragraphs, article.text, header.Title, header.image_id, header.date, header.views, header.shortDescription, image.desktop, image.mobile FROM post JOIN innerDescription ON post.innerAdvertising_id = innerDescription.innerDescription_id LEFT JOIN article ON post.post_id = article.post_id LEFT JOIN header ON post.post_id = header.post_id LEFT JOIN image ON post.post_id = image.image_id where post.post_id = $1 ORDER BY post.created_at DESC",
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
		&postModel.Im.Desktop,
		&postModel.Im.Mobile,
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

	rows, err := r.store.db.Query("SELECT post.post_id, innerDescription.innerAdvertising, post.created_at, article.Title, article.backgroundImg, article.paragraphs, article.text, header.Title, header.image_id, header.date, header.views, header.shortDescription, image.desktop, image.mobile FROM post JOIN innerDescription ON post.innerAdvertising_id = innerDescription.innerDescription_id LEFT JOIN article ON post.post_id = article.post_id LEFT JOIN header ON post.post_id = header.post_id LEFT JOIN image ON post.post_id = image.image_id ORDER BY post.created_at DESC")
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
			&postModel.Im.Desktop,
			&postModel.Im.Mobile,
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

	rows, err := r.store.db.Query("SELECT News_id, Title, Img, Date, Views, Created_at FROM news ORDER BY news.created_at DESC LIMIT 10")
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
