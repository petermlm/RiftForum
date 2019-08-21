package main

import (
    "fmt"
    "encoding/json"
    "net/http"

    "github.com/yalp/jsonpath"
)

func GetRedditHot(subreddit_name string) map[string]interface{} {
    subreddit_link := MakeRedditLink(subreddit_name)

    client := &http.Client{}
    req, err := http.NewRequest("GET", fmt.Sprintf("%s.json", subreddit_link), nil)
    req.Header.Set("User-Agent", "Riftforum")
    resp, err := client.Do(req)

    if err != nil {
    }

    defer resp.Body.Close()

    var results map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&results)
    return results
}

func MakeBBCodeListReddit(subreddit_name string, posts map[string]interface{}) string {
    subreddit_link := MakeRedditLink(subreddit_name)
    children_js, err := jsonpath.Read(posts, "$.data.children")

    if err != nil {
        return ""
    }

    list_str := fmt.Sprintf(
        "Currently in [url=%s]%s[/url]:\r\n[list]",
        subreddit_link,
        subreddit_name)


    switch children_obj := children_js.(type) {
    case interface{}:
        children := children_obj.([]interface{})

        slice_max := len(children)

        if slice_max == 0 {
            return ""
        }

        if slice_max > 5 {
            slice_max = 5
        }

        for i := range children[:slice_max] {
            child := children[i].(map[string]interface{})
            child_data := child["data"].(map[string]interface{})
            url := child_data["url"]
            title := child_data["title"]
            list_str += fmt.Sprintf("[*] [url=%s]%s[/url]\r\n", url, title)
        }
    default:
        return ""
    }

    list_str += "[/list]"
    return list_str
}
