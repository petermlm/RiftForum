package main

import "math/rand"

var banner_sentences = [...]string{
	"Now with 100% more Go",
	"Powered by PostgreSQL",
	"Powered by Redis",
	"Powered by Linux",
	"Powered by Markdown",
	"Malware Free",
	"Password == \"Banana7\"",
	"NPM Free",
	"Поддерживает UTF-8",
	"User Friendly",
	"No lollygagging",
	"No lurking",
}

func get_banner_sentence() string {
	return banner_sentences[rand.Intn(len(banner_sentences))]
}
