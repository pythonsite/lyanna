package models

type Tag struct {
	BaseModel
	Name string
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
