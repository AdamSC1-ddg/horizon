package hal

import (
	"net/url"
	"strings"
)

type LinkBuilder struct {
	Base *url.URL
}

func (lb *LinkBuilder) Link(parts ...string) Link {
	path := strings.Join(parts, "/")

	var href string
	if lb.Base != nil {
		pu, err := url.Parse(path)
		if err != nil {
			panic(err)
		}
		href = lb.Base.ResolveReference(pu).String()
	} else {
		href = path
	}

	return NewLink(href)
}

func (lb *LinkBuilder) PagedLink(parts ...string) Link {
	nl := lb.Link(parts...)
	nl.Href += StandardPagingOptions
	nl.PopulateTemplated()
	return nl
}
