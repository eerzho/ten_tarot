package entity

type TGMessage struct {
	ID     string `bson:"_id,omitempty" json:"id"`
	ChatID string `bson:"chat_id" json:"chat_id"`
	Text   string `bson:"text" json:"text"`
	Answer string `bson:"answer,omitempty" json:"answer"`
	File   string `bson:"file,omitempty" json:"file"`
}
