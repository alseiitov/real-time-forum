package model

type User struct {
	ID         int         `json:"id,omitempty"`
	Username   string      `json:"username,omitempty"`
	FirstName  string      `json:"firstName,omitempty"`
	LastName   string      `json:"lastName,omitempty"`
	Age        int         `json:"age,omitempty"`
	Gender     int         `json:"gender,omitempty"`
	Email      string      `json:"email,omitempty"`
	Password   string      `json:"-"`
	Registered interface{} `json:"registered,omitempty"`
	Role       int         `json:"role,omitempty"`
	Avatar     string      `json:"avatar,omitempty"`
}
