package models

type User struct {
	BaseModel
	Intro string
	Email string
	Name string `gorm:"unique_index"`
	PassWord string
	GitHubUrl string
	Active bool `gorm:"default:'1'"`
}

func(user *User) Insert() error {
	return 	DB.Create(user).Error
}

func(user *User) Update() error {
	return DB.Model(user).Updates(map[string]interface{}{
		"name": user.Name,
		"email":user.Email,
		"password": user.PassWord,
		"active":user.Active,
	}).Error
}

func (user *User) GetUserName(userID int) (string,error){
	err := DB.First(&user,userID).Error
	return user.Name,err
}

func GetUserByID(id interface{})(*User, error) {
	var user User
	err := DB.First(&user,id).Error
	return &user,err
}

func GetUserByName(username string)(*User, error) {
	var user User
	err := DB.First(&user, "Name=?",username).Error
	return &user,err
}

func ListUsers()([]*User, error) {
	var users []*User
	err := DB.Find(&users).Error
	return users,err
}

func GetUserNameByID(userID int)(name string,err error) {
	var user User
	err = DB.First(&user,"id=?",userID).Error
	return user.Name,err
}