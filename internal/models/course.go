package models

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Course struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Poster   string   `json:"poster"`
	Tags     []string `json:"tags"`
	Articles []int    `json:"articles"`
}

type Article struct {
	ID    int      `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	Text  string   `json:"text"`
}
