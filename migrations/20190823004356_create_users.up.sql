CREATE TABLE "images" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "desktop" varchar(100) DEFAULT NULL,
  "mobile" varchar(100) DEFAULT NULL
);
 
CREATE TABLE "news" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "title" varchar(100),
  "img" varchar(100),
  "date" datetime DEFAULT (now()),
  "views" int DEFAULT 0,
  "created_at" datetime DEFAULT (now())
);

CREATE TABLE "posts" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "headers_id" int,
  "article_id" int,
  "innerAdvertising_id" int,
  "created_at" datetime DEFAULT (now())
);

CREATE TABLE "headers" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "title" varchar(100),
  "image_id" int,
  "date" datetime DEFAULT (now()),
  "views" int DEFAULT 0,
  "shortDescription" varchar(100),
  "created_at" datetime DEFAULT (now())
);

CREATE TABLE "innerDescription" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "innerAdvertising" varchar(100)
);

CREATE TABLE "articles" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "title" varchar(100),
  "backgroundImg" varchar(100),
  "paragraphs" [text],
  "text" text,
  "recommendations_id" int DEFAULT NULL,
  "interaction_id" int DEFAULT NULL,
  "created_at" datetime DEFAULT (now())
);

CREATE TABLE "articles_images" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "articles_id" int,
  "images_id" int
);

CREATE TABLE "posts_products" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "products_id" int,
  "posts_id" int
);

CREATE TABLE "recommendations" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "title" varchar(100),
  "text" text
);

CREATE TABLE "interaction" (
  "id" SERIAL UNIQUE PRIMARY KEY,
  "title" varchar(100),
  "items" [text]
);

ALTER TABLE "innerDescription" ADD FOREIGN KEY ("id") REFERENCES "posts" ("innerAdvertising_id");

ALTER TABLE "images" ADD FOREIGN KEY ("id") REFERENCES "headers" ("image_id");

ALTER TABLE "articles" ADD FOREIGN KEY ("id") REFERENCES "articles_images" ("articles_id");

ALTER TABLE "images" ADD FOREIGN KEY ("id") REFERENCES "articles_images" ("images_id");

ALTER TABLE "articles" ADD FOREIGN KEY ("id") REFERENCES "posts_products" ("posts_id");

ALTER TABLE "news" ADD FOREIGN KEY ("id") REFERENCES "posts_products" ("posts_id");

ALTER TABLE "headers" ADD FOREIGN KEY ("id") REFERENCES "posts" ("headers_id");

ALTER TABLE "articles" ADD FOREIGN KEY ("id") REFERENCES "posts" ("article_id");

ALTER TABLE "recommendations" ADD FOREIGN KEY ("id") REFERENCES "articles" ("recommendations_id");

ALTER TABLE "interaction" ADD FOREIGN KEY ("id") REFERENCES "articles" ("interaction_id");

COMMENT ON COLUMN "news"."date" IS 'Дата публикации';

COMMENT ON COLUMN "news"."views" IS 'Кол-во просмотров';

COMMENT ON COLUMN "posts"."headers_id" IS 'Айдишник, указывающий на таблицу с краткой информацией о посте';

COMMENT ON COLUMN "posts"."innerAdvertising_id" IS 'Айдишник, указывающий на таблицу с внутренней рекламой сайта';

COMMENT ON COLUMN "headers"."title" IS 'Заголовок';

COMMENT ON COLUMN "headers"."image_id" IS 'Картинки для desktop и mobile версий';

COMMENT ON COLUMN "headers"."date" IS 'Дата публикации';

COMMENT ON COLUMN "headers"."views" IS 'Кол-во просмотров';

COMMENT ON COLUMN "headers"."shortDescription" IS 'Краткое описание для поста';

COMMENT ON COLUMN "innerDescription"."innerAdvertising" IS 'Внутренняя реклама сайта';

COMMENT ON COLUMN "articles"."title" IS 'Название поста';

COMMENT ON COLUMN "articles"."text" IS 'Текст статьи';

COMMENT ON COLUMN "articles"."recommendations_id" IS 'Айдишник, указывающий на таблицу с рекомендациями к употреблению';

COMMENT ON COLUMN "articles"."interaction_id" IS 'Айдишник, указывающий на таблицу взаимодействия вкусов пива и еды';

COMMENT ON COLUMN "articles_images"."articles_id" IS 'Айдишник, указывающий на таблицу статей';

COMMENT ON COLUMN "articles_images"."images_id" IS 'Айдишник, указывающий на таблицу картинок для разных версий сайта';

COMMENT ON COLUMN "posts_products"."products_id" IS 'Айдишник, указывающий на таблицу продуктов в другой бд';

COMMENT ON COLUMN "posts_products"."posts_id" IS 'Айдишник, указывающий на таблицу articles или news';

COMMENT ON COLUMN "interaction"."items" IS 'Перечисление сочетаний';