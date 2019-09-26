# graphql-news
This is example how to use grapqhql-go

## Setting up database
```
CREATE TABLE news (
	news_id serial NOT NULL,
	title varchar(100) NULL,
	body text NULL,
	author_id serial NOT NULL,
	CONSTRAINT news_pk PRIMARY KEY (news_id)
);

CREATE TABLE authors (
	"name" varchar(100) NULL,
	email varchar(100) NULL,
	author_id serial NOT NULL,
	CONSTRAINT author_pk PRIMARY KEY (author_id),
	CONSTRAINT author_un UNIQUE (email)
);
```