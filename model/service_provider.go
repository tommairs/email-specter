package model

import "regexp"

type ServiceProvider struct {
	Name          string         `json:"name"`
	Regex         string         `json:"regex"`
	CompiledRegex *regexp.Regexp `json:"-"`
}
