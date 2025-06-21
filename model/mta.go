package model

type MTA struct {
	ID            string `json:"id" bson:"_id,omitempty"`
	Name          string `json:"name" bson:"name"`
	SecretToken   string `json:"secret_token" bson:"secret_token"`
	CollectionUrl string `json:"collection_url,omitempty" bson:"-"`
}
