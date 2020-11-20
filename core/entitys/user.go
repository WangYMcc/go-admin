package entitys

type User struct {
	Id       int64 `orm:"column(id);unique;pk"`
	UserKey string `orm:"column(userkey);unique"`
	Username string `orm:"column(username);unique"`
	Password string `orm:"column(password);null"`
}

func (u User)TableName() string{
	return "user"
}
