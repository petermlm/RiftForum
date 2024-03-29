package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	redis "github.com/go-redis/redis"
)

var bots_initialized = false
var redis_client *redis.Client
var new_users_ch chan *User
var new_topics_ch chan *Topic
var new_messages_ch chan *Message

func InitBots() {
	// Initialize Redis Connection
	redis_client = redis.NewClient(&redis.Options{
		Addr:     Config.RedisAddr,
		Password: "",
		DB:       0,
	})

	_, err := redis_client.Ping().Result()
	if err != nil {
		panic("Can't connect with redis")
	}
	fmt.Println("Redis connection established")

	// Initialize channels
	new_users_ch = make(chan *User, Config.BotChannelLag)
	new_topics_ch = make(chan *Topic, Config.BotChannelLag)
	new_messages_ch = make(chan *Message, Config.BotChannelLag)

	// Fan out new messages
	new_msg_ch := make([]chan *Message, 2)
	for i, _ := range new_msg_ch {
		new_msg_ch[i] = make(chan *Message, Config.BotChannelLag)
	}

	go func() {
		for i := range new_messages_ch {
			for _, c := range new_msg_ch {
				c <- i
			}
		}
	}()

	// Stat bot functions
	go GreeterBot(new_users_ch)
	go RedditBot(new_msg_ch[0])
	go YoutubeBot(new_msg_ch[1])

	bots_initialized = true
	fmt.Println("Bots Started")
}

func SendNewUser(user *User) {
	if bots_initialized {
		new_users_ch <- user
	}
}

func SendNewTopic(topic *Topic) {
	if bots_initialized {
		new_topics_ch <- topic
	}
}

func SendNewMessage(message *Message) {
	if bots_initialized {
		new_messages_ch <- message
	}
}

func GetHearthBeatStatus() map[string]bool {
	ret := make(map[string]bool)

	key_pattern := redis_keys_pattern()
	keys, err := redis_client.Keys(key_pattern).Result()

	if err != nil {
		return ret
	}

	for i := range keys {
		key := keys[i]
		bot_name := get_bot_name_from_key(key)
		val, err := redis_client.Get(key).Result()

		if err != nil {
			ret[bot_name] = false
		} else {
			ret[bot_name] = hearthbeat_is_alive(val)
		}
	}

	return ret
}

func GreeterBot(new_users_ch chan *User) {
	greeter_template := `[center][b][color=red]Welcome to RiftForum, %s![/color][/b][/center]
    Thank you for Registering! You are welcome to [i][color=blue]post anything[/color][/i]. [b]BBCode[/b] can be used to style messages, and your signature can be edited in the [url=%s]user details page[/url]. User types are:
    [list]
    [*] Administrator
    [*] Moderator
    [*] Basic
    [*] Bot
    [/list]
    Happy Posting!`
	bot_name := "GreeterBot"
	user, _ := GetUser(bot_name)

	for {
		beat_hearth(bot_name)

		select {
		case new_user := <-new_users_ch:
			title := fmt.Sprintf("Welcome %s!", new_user.Username)
			user_detail_page := fmt.Sprintf("%s/users/%s", MakeBaseURL(), new_user.Username)
			message := fmt.Sprintf(greeter_template, new_user.Username, user_detail_page)
			NewTopic(user, title, message)
		case <-time.After(Config.BotHearthBeatPeriod * time.Second):
			func() {}()
		}
	}
}

func RedditBot(new_messages_ch chan *Message) {
	bot_name := "RedditBot"
	user, _ := GetUser(bot_name)
	re := regexp.MustCompile(`\/r\/[a-zA-Z0-9_]+`)

	for {
		beat_hearth(bot_name)

		select {
		case new_message := <-new_messages_ch:
			if new_message.Author.Id == user.Id {
				break
			}

			matches := re.FindAllString(new_message.Message, -1)

			if len(matches) == 0 {
				break
			}

			new_message_text := ""
			for i := range matches {
				posts := GetRedditHot(matches[i])
				bbcode_list := MakeBBCodeListReddit(matches[i], posts)

				if bbcode_list == "" {
					continue
				}

				new_message_text += bbcode_list
			}

			if new_message_text == "" {
				break
			}

			NewMessage(user, new_message.Topic, new_message_text)
		case <-time.After(Config.BotHearthBeatPeriod * time.Second):
			func() {}()
		}
	}
}

func YoutubeBot(new_messages_ch chan *Message) {
	bot_name := "YoutubeBot"
	user, _ := GetUser(bot_name)
	re := regexp.MustCompile(`^!youtubelist`)

	for {
		beat_hearth(bot_name)
		select {
		case new_message := <-new_messages_ch:
			if new_message.Author.Id == user.Id {
				break
			}

			if !re.MatchString(new_message.Message) {
				break
			}

			yt_list := youtube_videos_list(new_message.Topic)

			if yt_list == "" {
				break
			}

			NewMessage(user, new_message.Topic, yt_list)
		case <-time.After(Config.BotHearthBeatPeriod * time.Second):
			func() {}()
		}
	}
}

func ExpAnswerBot() {
	user, _ := GetUser("RiftBot")
	for {
		beat_hearth("ExpAnswerBot")

		select {
		case new_topic := <-new_topics_ch:
			NewMessage(user, new_topic, "Answer to a topic")
		case new_message := <-new_messages_ch:
			if new_message.Author.Id != user.Id {
				NewMessage(user, new_message.Topic, "Answer to a message")
			}
		case <-time.After(Config.BotHearthBeatPeriod * time.Second):
			func() {}()
		}
	}
}

func redis_key_name(bot_name string) string {
	return fmt.Sprintf("Hearthbeat_%s", bot_name)
}

func redis_keys_pattern() string {
	return "Hearthbeat_*"
}

func beat_hearth(bot_name string) {
	key_name := redis_key_name(bot_name)
	redis_client.Set(key_name, time.Now().Unix(), Config.BotHearthBeatExpire*time.Second)
}

func hearthbeat_is_alive(val string) bool {
	valint, err := strconv.ParseInt(val, 10, 64)

	if err != nil {
		return false
	}

	hearthbeat_interval := time.Now().Unix() - valint
	dead_interval := int64(Config.BotHearthBeatDead * time.Second)
	return hearthbeat_interval < dead_interval
}

func get_bot_name_from_key(key string) string {
	prefix_len := len(redis_keys_pattern())
	return key[prefix_len-1:]
}
