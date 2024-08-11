package model

import "bot/internal/constant"

type (
	TGInvoice struct {
		ID            string               `bson:"_id" json:"id"`
		Type          constant.InvoiceType `bson:"type" json:"type"`
		ChatID        string               `bson:"chat_id" json:"chat_id"`
		ChargeID      string               `bson:"charge_id" json:"charge_id,omitempty"`
		StarsCount    int                  `bson:"stars_count" json:"stars_count"`
		QuestionCount int                  `bson:"question_count" json:"question_count"`
		CreatedAt     string               `bson:"created_at" json:"created_at"`
	}
)
