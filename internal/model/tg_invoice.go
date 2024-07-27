package model

type (
	TGInvoice struct {
		ID            string `bson:"_id" json:"id"`
		ChatID        string `bson:"chat_id" json:"chat_id"`
		Stars         int    `bson:"stars" json:"stars"`
		QuestionCount int    `bson:"question_count" json:"question_count"`
		ChargeID      string `bson:"charge_id" json:"charge_id,omitempty"`
	}
)
