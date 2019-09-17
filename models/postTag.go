package models

import (
	"fmt"
	"log"
)

type PostTag struct {
	BaseModel
	PostID int64
	TagID int64
}

//func (pt *PostTag) Insert() error {
//	return DB.FirstOrCreate(pt, "post_id = ? and tag_id = ?", pt.PostId, pt.TagId).Error
//}

func InsertPostTag(postID, TagID int64) {
	postTag := PostTag{
		PostID:postID,
		TagID:TagID,
	}
	DB.Save(&postTag)
}

func GetPostsByTags(postID int64,tagids []int64)[]*Post {
	var posts []*Post
	_ = DB.Raw("select p.* from post_tags pt inner join posts p on p.id= pt.post_id where p.id != ? and pt.tag_id not in (?)",postID,tagids).Find(&posts).Error
	return posts
}

func ListTagByPostID (id interface{}) ([]*Tag,error) {
	var tags []*Tag
	rows,err := DB.Raw("select t.* from tags t inner join post_tags pt on t.id = pt.tag_id where pt.post_id = ?",id).Rows()
	if err != nil {
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Tag
		_ = DB.ScanRows(rows, &tag)
		tags = append(tags, &tag)
	}
	return tags, nil
}

func GetTagNames (tags []*Tag) []string {
	var tagNames []string
	for _,v :=range tags {
		tagNames = append(tagNames,v.Name)
	}
	return tagNames
}

func DeleteTagByPostID(postID interface{}) error {
	return DB.Delete(&PostTag{},"post_id=?",postID).Error
}

func UpdateMultiTags(originTags []string, newTags []string, postID int) {
	needToDelTags := GetTagArray(originTags, newTags)
	var needToDelTagID []int
	for _,v := range needToDelTags {
		tagID := GetTagIDByName(v)
		needToDelTagID = append(needToDelTagID, tagID)
	}
	DB.Delete(&PostTag{},"tag_id in ( ? )",needToDelTagID)

	needToAddTags := GetTagArray(newTags, originTags)
	fmt.Println(needToAddTags)
	var needAddTagID []int
	for _,v := range needToAddTags {
		tag := Tag{
			Name:v,
		}
		GetTag(&tag)
		needAddTagID = append(needAddTagID,int(tag.ID))
		log.Println(needToDelTagID)
	}
	for _,v := range needAddTagID {
		InsertPostTag(int64(postID),int64(v))
	}




}

// 判断tag是否在array中
func IsInArray(name string, tagNameArray []string) bool {
	for _, v:= range tagNameArray {
		if name == v {
			return  true
		}
	}
	return false
}

// 获取需要删除的tags以及需要新增加的tags
func GetTagArray(originTags []string, newTags []string) []string {
	var tags []string
	for _,v1 := range originTags {
		if !IsInArray(v1,newTags) {
			tags = append(tags, v1)
		}
	}
	return tags
}