package graphserver

type Author struct {
	AuthorID int    `json:"author_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func (gs *GraphServer) GetAuthorByID(authorID int) (author *Author, err error) {
	author = &Author{}
	err = gs.DB.QueryRow("select author_id, name, email from authors where author_id = $1", authorID).Scan(&author.AuthorID, &author.Name, &author.Email)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (gs *GraphServer) InsertAuthor(author *Author) (authorID int, err error) {
	var lastInsertID int
	err = gs.DB.QueryRow("INSERT INTO authors(name, email) VALUES($1, $2) returning author_id;", author.Name, author.Email).Scan(&lastInsertID)
	if err != nil {
		return lastInsertID, err
	}
	return lastInsertID, nil
}

func (gs *GraphServer) UpdateAuthor(authorID int, author *Author) (authorReturn *Author, err error) {
	stmt, err := gs.DB.Prepare("UPDATE authors SET name = $1, email = $2 WHERE author_id = $3")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(author.Name, author.Email, authorID)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (gs *GraphServer) GetAllAuthor() (authors []*Author, err error) {

	rows, err := gs.DB.Query("SELECT author_id, name, email FROM authors")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		author := &Author{}

		err = rows.Scan(&author.AuthorID, &author.Name, &author.Email)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}
