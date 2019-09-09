package models


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

func DeleteTagByPostID(postID interface{}) error {
	return DB.Delete(&PostTag{},"post_id=?",postID).Error
}