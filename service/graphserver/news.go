package graphserver

import "log"

type News struct {
	NewsID   int    `json:"news_id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (gs *GraphServer) GetNewsByID(newsID int) (news *News, err error) {
	news = &News{}
	err = gs.DB.QueryRow("select id, title, body, author_id from news where id = $1", newsID).Scan(&news.NewsID, &news.Title, &news.Body, &news.AuthorID)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (gs *GraphServer) InsertNews(news *News) (newsID int, err error) {
	var lastInsertID int
	err = gs.DB.QueryRow("INSERT INTO news(title, body, author_id) VALUES($1, $2, $3) returning news_id;", news.Title, news.Body, news.AuthorID).Scan(&lastInsertID)
	if err != nil {
		return lastInsertID, err
	}
	return lastInsertID, nil
}

func (gs *GraphServer) UpdateNews(newsID int, news *News) (newsReturn *News, err error) {
	stmt, err := gs.DB.Prepare("UPDATE news SET title = $1, body = $2, author_id = $3 WHERE id = $3")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(news.Title, news.Body, news.AuthorID, newsID)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (gs *GraphServer) GetAllNews() (arrNews []*News, err error) {

	rows, err := gs.DB.Query("SELECT id, title, body, author_id FROM news")
	log.Println(err)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		news := &News{}

		err = rows.Scan(&news.NewsID, &news.Title, &news.Body, &news.AuthorID)
		if err != nil {
			return nil, err
		}

		arrNews = append(arrNews, news)
	}
	return arrNews, nil
}
