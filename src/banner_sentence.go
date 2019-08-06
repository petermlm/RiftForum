package main

import "math/rand"

var banner_sentences = [...]string{
    "Now with 100% more Go",
    "Malware Free",
    "Powered by Markdown",
    "Password == \"Banana7\"",
    "All your base are belong to us",
    "Du Hast",
    "Spaces > Tabs",
    "Minecraft Reference",
    "No JS",
    "Irony",
    "Підтримує UTF-8",
}

func set_banner_sentence() string {
    return banner_sentences[rand.Intn(len(banner_sentences))]
}
