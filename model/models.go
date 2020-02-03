package model

type Repository struct {
	Id   uint64 `json:"id,omitempty"`
	User string `json:"user,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
	Desc string `json:"description,omitempty"`
	Lang string `json:"language,omitempty"`
	Tags []Tag  `json:"tags,omitempty"`
}

type Tag struct {
	Name string `json:"name"`
}

type TagUpdate struct {
	OldTag Tag `json:"oldtag"`
	NewTag Tag `json:"newtag"`
}
