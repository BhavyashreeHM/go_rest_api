package models

type Student struct {
	Id        int `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"lastt_name,omitempty"`
	Class     string `json:"class,omitempty"`
	Subject   string `json:"subject,omitempty"`
}
