package models

type Post struct {
	BaseModel
	User User `gorm:"-"`
	Title string
	AuthorID int
	Slug string
	Summary string
	CanComment bool
	Published bool
	Tags []*Tag `gorm:"-"`
}

func(post *Post) Insert() error {
	return 	DB.Create(post).Error
}

func (post *Post) Update() {
	DB.Save(post)
}

func (post *Post) GetUserName(userID int)string {
	Name, _ := post.User.GetUserName(userID)
	return Name
}

func ListPosts()([]*Post, error) {
	var posts []*Post
	err := DB.Find(&posts).Error
	return posts,err
}

func GetPostByID(postID interface{})(*Post,error) {
	var post Post
	err := DB.First(&post,postID).Error
	return &post,err
}

func (post *Post) IsInAllTags(tagName string, tagNames []string) bool {
	for _, v := range tagNames {
		if v == tagName {
			return true
		}
	}
	return false
}


