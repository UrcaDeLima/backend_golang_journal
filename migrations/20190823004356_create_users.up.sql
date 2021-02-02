CREATE TABLE users (
  id bigserial not null primary key,
  email varchar not null unique,
  encrypted_password varchar not null
);







CREATE TABLE image (
  image_id SERIAL PRIMARY KEY,
  desktop varchar(255) DEFAULT NULL,
  mobile varchar(255) DEFAULT NULL
);
 
CREATE TABLE news (
  news_id SERIAL PRIMARY KEY,
  title varchar(255),
  img varchar(255),
  date timestamp NOT NULL DEFAULT NOW(),
  views INTEGER DEFAULT 0,
  created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE innerDescription (
  innerDescription_id SERIAL PRIMARY KEY,
  innerAdvertising varchar(255)
);

CREATE TABLE post (
  post_id SERIAL PRIMARY KEY,
  innerAdvertising_id INTEGER REFERENCES innerDescription (innerDescription_id) ON DELETE CASCADE,
  created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE header (
  header_id SERIAL PRIMARY KEY,
  title varchar(255),
  image_id INTEGER,
  date timestamp NOT NULL DEFAULT NOW(),
  views int DEFAULT 0,
  shortDescription varchar(255),
  post_id INTEGER REFERENCES post (post_id) ON DELETE CASCADE,
  created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE article (
  article_id SERIAL PRIMARY KEY,
  title varchar(255),
  backgroundImg varchar(255),
  paragraphs text[],
  text text,
  post_id INTEGER REFERENCES post (post_id) ON DELETE CASCADE,
  created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE article_image (
  article_id INTEGER REFERENCES article (article_id),
  image_id INTEGER REFERENCES image (image_id)
);

CREATE TABLE article_product (
  product_id INTEGER,
  article_id INTEGER REFERENCES article (article_id) ON DELETE CASCADE
);

CREATE TABLE recommendation (
  recommendation_id SERIAL PRIMARY KEY,
  title varchar(255),
  article_id INTEGER REFERENCES article (article_id) ON DELETE CASCADE,
  text text
);

CREATE TABLE interaction (
  interaction_id SERIAL PRIMARY KEY,
  title varchar(255),
  article_id INTEGER REFERENCES article (article_id) ON DELETE CASCADE,
  items text[]
);

COMMENT ON COLUMN "news"."date" IS 'Дата публикации';

COMMENT ON COLUMN "news"."views" IS 'Кол-во просмотров';

COMMENT ON COLUMN "header"."title" IS 'Заголовок';

COMMENT ON COLUMN "header"."image_id" IS 'Картинки для desktop и mobile версий';

COMMENT ON COLUMN "header"."date" IS 'Дата публикации';

COMMENT ON COLUMN "header"."views" IS 'Кол-во просмотров';

COMMENT ON COLUMN "article"."title" IS 'Название поста';

COMMENT ON COLUMN "article"."text" IS 'Текст статьи';

COMMENT ON COLUMN "article_image"."article_id" IS 'Айдишник, указывающий на таблицу статей';

COMMENT ON COLUMN "article_image"."image_id" IS 'Айдишник, указывающий на таблицу картинок для разных версий сайта';

COMMENT ON COLUMN "article_product"."product_id" IS 'Айдишник, указывающий на таблицу продуктов в другой бд';

COMMENT ON COLUMN "article_product"."article_id" IS 'Айдишник, указывающий на таблицу articles';

COMMENT ON COLUMN "interaction"."items" IS 'Перечисление сочетаний';