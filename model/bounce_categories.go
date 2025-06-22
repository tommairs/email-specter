package model

import (
	"github.com/dlclark/regexp2"
)

type BounceCategories struct {
	Blacklist []regexp2.Regexp `json:"blacklist"`
	Greylist  []regexp2.Regexp `json:"greylist"`
	Spam      []regexp2.Regexp `json:"spam"`
	RateLimit []regexp2.Regexp `json:"rate_limit"`
	Block     []regexp2.Regexp `json:"block"`
	Virus     []regexp2.Regexp `json:"virus"`
	Recipient []regexp2.Regexp `json:"recipient"`
	Message   []regexp2.Regexp `json:"message"`
	Config    []regexp2.Regexp `json:"config"`
	Sender    []regexp2.Regexp `json:"sender"`
	Dmarc     []regexp2.Regexp `json:"dmarc"`
	Policy    []regexp2.Regexp `json:"policy"`
	Capacity  []regexp2.Regexp `json:"capacity"`
	Envelope  []regexp2.Regexp `json:"envelope"`
	Network   []regexp2.Regexp `json:"network"`
	Protocol  []regexp2.Regexp `json:"protocol"`
	Auth      []regexp2.Regexp `json:"auth"`
	Other     []regexp2.Regexp `json:"other"`
}

type CategoryEntry struct {
	Name     string
	Patterns []regexp2.Regexp
}

func (b *BounceCategories) GetCategories() []CategoryEntry {

	return []CategoryEntry{
		{"Blacklist", b.Blacklist},
		{"Greylist", b.Greylist},
		{"Spam", b.Spam},
		{"RateLimit", b.RateLimit},
		{"Block", b.Block},
		{"Virus", b.Virus},
		{"Recipient", b.Recipient},
		{"Message", b.Message},
		{"Config", b.Config},
		{"Sender", b.Sender},
		{"Dmarc", b.Dmarc},
		{"Policy", b.Policy},
		{"Capacity", b.Capacity},
		{"Envelope", b.Envelope},
		{"Network", b.Network},
		{"Protocol", b.Protocol},
		{"Auth", b.Auth},
		{"Other", b.Other},
	}

}
