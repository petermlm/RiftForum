package main

import (
    "strconv"
    "net/http"
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

func PageDefault() Page {
    return Page{
        limit: PageDefaultLimit,
        offset: PageDefaultOffset,
    }
}

func PageFromRequest(req *http.Request) Page {
    page_size := req.URL.Query().Get("page_size")
    page_num := req.URL.Query().Get("page_num")

    page_size_parsed, page_size_err := strconv.ParseInt(page_size, 10, 32)
    page_num_parsed, page_num_err := strconv.ParseInt(page_num, 10, 32)

    var page Page

    if page_size_err != nil || page_size_parsed <= 0 {
        page_size_parsed = 20
    }

    if page_num_err != nil {
        page_num_parsed = 0
    }

    page.limit = int(page_size_parsed)
    page.offset = int(page_size_parsed * page_num_parsed)

    return page
}
