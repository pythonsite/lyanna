package models

type Tag struct {
	BaseModel
	Name string
	Total int `gorm:"-"`
}

func ListALlTags()([]*Tag, error) {
	var tags []*Tag
	err := DB.Find(&tags).Error
	return tags,err
}

func GetTagIDByName(name string) int {
	var tag Tag
	DB.First(&tag,"name=?",name)
	return int(tag.ID)
}

func GetTag(tag *Tag){
	res := DB.FirstOrCreate(tag,"name=?",tag.Name).Row()
	_ = res.Scan(tag)
}

func GetTagNameByID(tagID int)string {
	var tag Tag
	DB.First(&tag,"id=?",tagID)
	return tag.Name
}

func ListTag()([]*Tag,error) {
	var tags []*Tag
	rows, err := DB.Raw("select t.*,count(*) total from tags t inner join post_tags pt on t.id=pt.tag_id inner join posts p on pt.post_id = p.id where p.published = ? group by pt.tag_id",true).Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tag Tag
		DB.ScanRows(rows,&tag)
		tags = append(tags, &tag)
	}
	return tags, nil
}
