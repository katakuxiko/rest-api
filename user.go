package todo

type User struct {
	id int `json:"-"`
	Name string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}