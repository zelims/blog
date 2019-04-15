package models

type User struct {
	ID			int
	Username	string
	Password	string
}

type UserProfile struct {
	Name		string
	Location	string
	About		string
	Github		string
	Twitter		string
	Instagram	string
	LinkedIn	string
}