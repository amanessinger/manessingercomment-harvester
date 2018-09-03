package harvest

type Comment struct {
	Id          int    `json:"id"`
	Date        string `json:"date"`
	AttachTo    string `json:"attach_to"`
	Author      string `json:"author"`
	AuthorEmail string `json:"author_email"`
	AuthorUrl   string `json:"author_url"`
	Content     string `json:"content"`
	IndentLevel int    `json:"indent_level"`
}
