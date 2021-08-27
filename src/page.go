package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type Page struct {
	limit, offset int
}

func (p Page) get_size() int {
	return p.limit
}

func (p Page) get_num() int {
	return p.offset / p.limit
}

func (p Page) querystr() string {
	return fmt.Sprintf("page_size=%d&page_num=%d", p.get_size(), p.get_num())
}

func NewPage(num int) Page {
	return Page{
		limit:  Config.PageDefaultSize,
		offset: num * Config.PageDefaultLimit,
	}
}

func PageDefault() Page {
	return Page{
		limit:  Config.PageDefaultLimit,
		offset: Config.PageDefaultOffset,
	}
}

func PageFromRequest(req *http.Request) Page {
	page_size := req.URL.Query().Get("page_size")
	page_num := req.URL.Query().Get("page_num")

	page_size_parsed_64, page_size_err := strconv.ParseInt(page_size, 10, 64)
	page_num_parsed_64, page_num_err := strconv.ParseInt(page_num, 10, 64)

	page_size_parsed := int(page_size_parsed_64)
	page_num_parsed := int(page_num_parsed_64)

	var page Page

	if page_size_err != nil || page_size_parsed <= 0 {
		page_size_parsed = Config.PageDefaultSize
	}

	if page_num_err != nil {
		page_num_parsed = Config.PageDefaultNum
	}

	page.limit = int(page_size_parsed)
	page.offset = int(page_size_parsed * page_num_parsed)

	return page
}
