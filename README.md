# graphql-news
This is example how to crud using grapqhql-go 
- graphql-go
- postgresql

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

## Update dependencies
```
go get github.com/graphql-go/graphql
go get github.com/graphql-go/handler
go get github.com/lib/pq
```

## usage
Access `localhost:8081/graphql`

Query to get the all authors
```
query {
  authors {
    author_id
    name
    email
  }
}
```

Query to get single author
```
query {
  author(author_id:1) {
    author_id
    name
    email
  }
}
```

Query to insert news
```
mutation {
  insert_news(body: "ini tulisan dalemnya", title: "ini titlenya cuy", author_id:1) {
    body
    author {
      author_id
      email
    }
  }
}

```