package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
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

func GetPosts(postID int64) {
	tags, _ := ListTagByPostID(postID)
	var tagids []int64
	for _,tag := range tags {
		tagids = append(tagids,int64(tag.ID))
	}
	posts := GetPostsByTags(postID,tagids)

}

type Archive struct {
	ArchiveDate time.Time
	Total int
	Year string
}

func ListPostArchives()([]*Archive, error) {
	var archives []*Archive
	rows, _ := DB.Raw("select DATE_FORMAT(created_at,'%Y') as year, count(*) as total from posts where published = ? group by year order by year desc",true).Rows()
	defer rows.Close()
	for rows.Next() {
		var archive Archive
		DB.ScanRows(rows, &archive)
		archive.ArchiveDate, _ =  time.Parse("2006", archive.Year)
		archives = append(archives, &archive)
	}
	return archives, nil
}

func ListPostByArchive(year string)[]*Post {
	condition := fmt.Sprintf("%s",year)
	rows, _ := DB.Raw("select * from posts where date_format(created_at,'%Y')=? and published = ? order by created_at desc",condition,true).Rows()
	defer rows.Close()
	posts := make([]*Post,0)
	for rows.Next() {
		var post Post
		DB.ScanRows(rows,&post)
		posts = append(posts, &post)
	}
	return posts
}
