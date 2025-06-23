package model

import "github.com/dlclark/regexp2"

type BounceTypeEntry struct {
	CompiledRegex *regexp2.Regexp `json:"-"`
	Category      string          `json:"category"`
}
