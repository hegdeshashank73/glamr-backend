package utils

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
)

type StringOpts struct {
	EntityName         string
	ConvertToLowercase bool
	Slugify            bool
	ConvertToUppercase bool
	MaxLength          int
	IsRequired         bool
	IsMarkdown         bool
	EscapeWhitespace   bool
	IsLink             bool
}

func AssertString(txt string, opts StringOpts) (string, error) {
	txt = strings.TrimSpace(txt)
	if opts.MaxLength > 0 && len(txt) > opts.MaxLength {
		return "", fmt.Errorf("'%s' too long; max length could be %d", opts.EntityName, opts.MaxLength)
	}

	if opts.IsRequired && len(txt) == 0 {
		return "", fmt.Errorf("'%s' is required", opts.EntityName)
	}

	if opts.ConvertToLowercase {
		txt = strings.ToLower(txt)
	}
	if opts.ConvertToUppercase {
		txt = strings.ToUpper(txt)
	}
	if opts.Slugify {
		txt = slug.Make(txt)
	}
	if opts.EscapeWhitespace {
		txt = strings.ReplaceAll(txt, "\n", "\\n")
	}
	if opts.IsLink && len(txt) > 0 {
		if !strings.HasPrefix(txt, "http://") && !strings.HasPrefix(txt, "https://") {
			return "", fmt.Errorf("'%s' should be a valid link", opts.EntityName)
		}
	}
	return txt, nil
}

func StringOrDefault(args ...string) string {
	for _, a := range args {
		if a != "" {
			return a
		}
	}
	return ""
}
