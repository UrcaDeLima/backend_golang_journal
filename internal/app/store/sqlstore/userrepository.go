package sqlstore

import (
	"database/sql"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/UrcaDeLima/backend_golang_journal/internal/app/model"
	"github.com/UrcaDeLima/backend_golang_journal/internal/app/store"
	"github.com/disintegration/imaging"
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

// InteractionRepository ...
type InteractionRepository struct {
	store *Store
}

// RecommendationRepository ...
type RecommendationRepository struct {
	store *Store
}

// PostRepository ...
type PostRepository struct {
	store *Store
}

// createFile ...
func createFile(part *multipart.Part, image *model.Image) *model.Image {
	dst, err := os.Create("./image/" + part.FileName())
	if err != nil {
		fmt.Println(err)
	}

	if part.FormName() == "desktop" {
		image.Desktop = "image/" + part.FileName()
	} else if part.FormName() == "mobile" {
		image.Mobile = "image/" + part.FileName()
	}

	io.Copy(dst, part)

	file, err := os.Open("image/" + part.FileName())
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fileSize, err := jpeg.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", "image/"+part.FileName(), err)
	}

	if fileSize.Width > 1920 {
		// load original image
		img, err := imaging.Open("image/" + part.FileName())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		dstimg := imaging.Resize(img, 1920, 0, imaging.Box)

		// save resized image
		err = imaging.Save(dstimg, "image/"+part.FileName())

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return image
}

// UpdatePicture ...
func (r *ImageRepository) UpdatePicture(id int, m *multipart.Reader) error {
	image := &model.Image{}
	r.store.db.QueryRow(
		"SELECT Desktop, Mobile FROM image WHERE image_id = $1",
		id,
	).Scan(&image.Desktop, &image.Mobile)

	var directory string

	for {
		part, err := m.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			continue
		}

		if part.FormName() == "desktop" {
			directory = image.Desktop
		} else if part.FormName() == "mobile" {
			directory = image.Mobile
		}

		err = os.Remove(directory)
		if err != nil {
			fmt.Println(err)
		}

		image = createFile(part, image)
	}
	return r.store.db.QueryRow(
		"UPDATE image SET desktop = $1, mobile = $2 WHERE image_id = $3 RETURNING image_id",
		image.Desktop,
		image.Mobile,
		id,
	).Scan(&image.Image_id)
}

// SetPicture ...
func (r *ImageRepository) SetPicture(m *multipart.Reader) error {
	image := &model.Image{}

	for {
		part, err := m.NextPart()
		if err == io.EOF {
			break
		}
		log.Println(part)

		if part.FileName() == "" {
			continue
		}

		image = createFile(part, image)
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

	log.Println(rows)
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
		log.Println(postModel)
		allPostModel = append(allPostModel, postModel)
	}
	if err := rows.Err(); err != nil {
		log.Println(allPostModel)
		log.Fatal(err)
	}

	return allPostModel, err
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func dbQueryCreatePost(postModel *model.PostModel) {
	db, err := sql.Open(os.Getenv("DB_DIALECT"), os.Getenv("DB_URL"))
	handleError(err)
	tx, err := db.Begin()
	handleError(err)
	defer db.Close()
	var post_id int

	if err := tx.QueryRow(
		"INSERT INTO innerDescription (innerAdvertising) VALUES ($1) RETURNING innerDescription_id",
		&postModel.I.InnerAdvertising,
	).Scan(
		&post_id,
	); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}

	//log.Println(post_id)
	if err := tx.QueryRow(
		"INSERT INTO post (innerAdvertising_id) VALUES ($1) RETURNING post_id",
		post_id,
	).Scan(
		&post_id,
	); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	//log.Println(post_id)

	var imageID int
	if err := tx.QueryRow(
		"INSERT INTO image (desktop, mobile) VALUES ($1, $2) RETURNING image_id",
		&postModel.Im.Desktop,
		&postModel.Im.Mobile,
	).Scan(
		&imageID,
	); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	log.Println(postModel.H.Title)
	log.Println(postModel.H.Views)
	log.Println(postModel.H.ShortDescription)
	log.Println(imageID)

	// if err := tx.QueryRow(
	// 	"INSERT INTO header (title, image_id, views, shortDescription) VALUES ($1, $2, $3, $4) RETURNING header_id", // понять что не так с запросом
	// 	postModel.H.Title,
	// 	imageID,
	// 	postModel.H.Views,
	// 	postModel.H.ShortDescription,
	// ).Scan(
	// 	&post_id,
	// ); err != nil {
	// 	log.Println("here")
	// 	tx.Rollback()
	// 	log.Fatal(err)
	// }

	if err := tx.QueryRow(
		"INSERT INTO header (title, image_id, views, shortDescription) VALUES ('test', 2, 3, 'test') RETURNING header_id", // понять что не так с запросом
	).Scan(
		&post_id,
	); err != nil {
		log.Println("qwqwqw")
		tx.Rollback()
		log.Fatal(err)
	}

	log.Println(111111111)
	log.Println(imageID)

	var articleID int
	if err := tx.QueryRow(
		"INSERT INTO article (title, BackgroundImg, Paragraphs, Text, Post_id) VALUES ($1, $2, $3, $4, 19) RETURNING article_id",
		&postModel.A.Title,
		imageID,
		&postModel.A.Paragraphs,
		&postModel.A.Text,
	).Scan(
		&articleID,
	); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	log.Println(articleID)

	if err := tx.QueryRow(
		"INSERT INTO article_image (Article_id, image_id) VALUES ($1, $2)",
		articleID,
		imageID,
	); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	log.Println(articleID)

	if err := tx.QueryRow(
		"INSERT INTO article_product (Article_id, product_id) VALUES ($1, 11)", // поставить динамический id
		articleID,
	); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	log.Println(articleID)

	if err := tx.QueryRow(
		"INSERT INTO recomendation (title, Article_id, Text) VALUES ($1, $2, $3)",
		&postModel.R.Title,
		articleID,
		&postModel.R.Text,
	); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	log.Println(articleID)

	if err := tx.QueryRow(
		"INSERT INTO interaction (title, Article_id, Items) VALUES ($1, $2, $3)",
		&postModel.In.Title,
		articleID,
		&postModel.In.Items,
	); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	log.Println(articleID)

	log.Println("all okey")

	tx.Commit()
}

// CreatePost ...
func (r *PostRepository) CreatePost(m *multipart.Reader) {
	postModel := &model.PostModel{}
	image := &model.Image{}
	//text := make([]byte, 512)
	//number := make([]byte, 512)

	for {
		part, err := m.NextPart()
		if err == io.EOF {
			break
		}

		//log.Println(part)
		if part.FileName() == "" {
			switch part.FormName() {
			case "headerShortDescription", "innerDescriptionInnerAdvertising", "interactionTitle", "recommendationTitle", "headerTitle", "articleTitle", "articleParagraphs", "articleText", "recommendationText", "interactionItems":
				{
					text := make([]byte, 512)
					_, err = part.Read(text)
					if err != nil && err != io.EOF {
						fmt.Println(err)
						return
					}
					log.Println(string(text))
					log.Print("\n")

					switch part.FormName() {
					case "headerShortDescription":
						{
							postModel.H.ShortDescription = string(text)
						}
					case "innerAdvertising":
						{
							postModel.I.InnerAdvertising = string(text)
						}
					case "interactionTitle":
						{
							postModel.In.Title = string(text)
						}
					case "recommendationTitle":
						{
							postModel.R.Title = string(text)
						}
					case "headerTitle":
						{
							postModel.H.Title = string(text)
						}
					case "articleTitle":
						{
							postModel.A.Title = string(text)
						}
					case "articleParagraphs":
						{

							//log.Println(string(text))
							//log.Print("\n")
							postModel.A.Paragraphs = append(postModel.A.Paragraphs, string(text))
						}
					case "articleText":
						{
							postModel.A.Text = string(text)
						}
					case "recommendationText":
						{
							postModel.R.Text = string(text)
						}
					case "interactionItems":
						{
							//log.Println(string(text))
							//log.Print("\n")
							postModel.In.Items = string(text)
						}
					}
				}
			case "headerViews":
				{
					// _, err = part.Read(number)
					// if err != nil && err != io.EOF {
					// 	fmt.Println(err)
					// 	return
					// }
					// views, _ := strconv.Atoi(string(text))
					// postModel.H.Views = int(views) // сделать парсинг в int, сейчас не работает
					//log.Println(text)
				}
			default:
				{
					log.Println("Error " + part.FormName())
				}
			}

			//log.Println(postModel)
			//continue
		} else {
			image = createFile(part, image)
			postModel.Im = *image
		}
	}
	dbQueryCreatePost(postModel)
}

// // CreatePost ...
// func (r *PostRepository) CreatePost(post *model.Post) error {
// 	var db *sql.DB

// 	return r.store.db.QueryRow(
// 		"INSERT INTO post (innerAdvertising_id) VALUES ($1) RETURNING post_id",
// 		post.InnerAdvertising_id,
// 	).Scan(&post.Post_id)
// }

// // CreateInnerDescription ...
// func (r *InnerDescriptionRepository) CreateInnerDescription(innerDescription *model.InnerDescription) error {
// 	return r.store.db.QueryRow(
// 		"INSERT INTO innerDescription (innerAdvertising) VALUES ($1) RETURNING innerDescription_id",
// 		innerDescription.InnerAdvertising,
// 	).Scan(&innerDescription.InnerDescription_id)
// }

// // CreateHeader ...
// func (r *HeaderRepository) CreateHeader(header *model.Header) error {
// 	return r.store.db.QueryRow(
// 		"INSERT INTO header (title, image_id, date, views, shortDescription, post_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING Header_id",
// 		header.Title,
// 		header.Image_id,
// 		header.Date,
// 		header.Views,
// 		header.ShortDescription,
// 		header.Post_id,
// 	).Scan(&header.Header_id)
// }

// // CreateArticle ...
// func (r *ArticleRepository) CreateArticle(article *model.Article) error {
// 	return r.store.db.QueryRow(
// 		"INSERT INTO header (title, backgroundImg, text, views, post_id) VALUES ($1, $2, $3, $4, $5) RETURNING article_id",
// 		article.Title,
// 		article.BackgroundImg,
// 		article.Paragraphs,
// 		article.Text,
// 		article.Post_id,
// 	).Scan(&article.Article_id)
// }

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
