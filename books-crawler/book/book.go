package book

type Book struct {
	Thumbnail string `json:"thumbnail"`
	DetailURL string `json:"detailURL"`
	Title     string `json:"title"`
	Rating    int    `json:"rating"`
	Price     string `json:"price"`
	Instock   bool   `json:"instock"`
}
