package models

import (
	"database/sql"
	"strconv"
)

var RedisPostKey string = "posts/%d/props/content"

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

func _listPost(tag string, published bool) ([]*Post, error){
	var posts  []*Post
	var err error
	if len(tag) > 0 {
		tagID, err := strconv.ParseUint(tag,10,64)
		if err != nil {
			return nil, err
		}
		var rows *sql.Rows
		if published {
			rows, err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id=? and p.published = ? order by created_at desc",tagID,true).Rows()
		} else {
			rows , err = DB.Raw("select p.* from posts p inner join post_tags pt on p.id=pt.post_id where pt.tag_id=? order by created_at desc",tagID).Rows()
		}
		if err != nil {
			return nil,err
		}
		defer rows.Close()
		for rows.Next(){
			var post Post
			DB.ScanRows(rows, &post)
			posts = append(posts, &post)
		}
	} else {
		if published {
			err = DB.Where("published = ?", true).Order("created_at desc").Find(&posts).Error
		} else {
			err = DB.Order("created_at desc").Find(&posts).Error
		}
	}
	return posts,err
}

func ListPublishedPost(tag string)([]*Post, error) {
	return _listPost(tag, true)
}

func CountPostByTag(tag string)(count int, err error) {
	var tagID uint64
	if len(tag) > 0 {
		tagID, err = strconv.ParseUint(tag,10,64)
		if err != nil {
			return
		}
		err = DB.Raw("select count(*) from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id=? and p.published=?",tagID,true).Row().Scan(&count)
	} else {
		err = DB.Raw("select count(*) from posts p where p where p.published=?",true).Row().Scan(&count)
	}
	return
}

func PostCreatAndGetID(post *Post)error {
	err := DB.Create(post).Row().Scan(post)
	return err
}

