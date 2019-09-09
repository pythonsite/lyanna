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
