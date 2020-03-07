package main

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

func youtube_videos_list(topic *Topic) string {
	yt_links := get_youtube_links_from_topic(topic)
	return make_bbcode_list_youtube(yt_links)
}

func get_youtube_links_from_topic(topic *Topic) map[string]bool {
	yt_links := make(map[string]bool)
	re := regexp.MustCompile(`https?://(www\.)?youtube\.com\/watch\?([\w-]*=[\w-]*)((&([\w-]*=[\w-]*))*)?`)

	for i := range topic.Messages {
		msg := topic.Messages[i]
		matches := re.FindAllString(msg.Message, -1)

		for i := range matches {
			match := matches[i]
			yt_links[match] = true
		}
	}

	return yt_links
}

func parse_yt_link(link string) (string, error) {
	yt_link := "https://www.youtube-nocookie.com/embed/"

	url_parsed, err := url.Parse(link)
	if err != nil {
		return "", errors.New("Not a url")
	}

	query, _ := url.ParseQuery(url_parsed.RawQuery)

	video_id_arr, ok := query["v"]
	if !ok {
		return "", errors.New("No video id")
	}

	video_id := video_id_arr[0]
	return yt_link + video_id, nil
}

func make_bbcode_list_youtube(yt_links map[string]bool) string {
	if len(yt_links) == 0 {
		return ""
	}

	list_str := fmt.Sprintf("Youtube links in this threat:\r\n[list]")

	for k := range yt_links {
		list_str += fmt.Sprintf("[*] [url]%s[/url]\r\n", k)
	}

	list_str += "[/list]"
	return list_str
}
