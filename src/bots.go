package main

import (
    "fmt"
    "time"
    "strconv"
    "regexp"

    redis "github.com/go-redis/redis"
)

var bots_initialized = false
var redis_client *redis.Client
var new_users_ch chan *User
var new_topics_ch chan *Topic
var new_messages_ch chan *Message

func InitBots() {
    redis_client = redis.NewClient(&redis.Options{
        Addr:     RedisAddr,
        Password: "",
        DB:       0,
    })

    _, err := redis_client.Ping().Result()
    if err != nil {
        panic("Can't connect with redis")
    }

    fmt.Println("Redis connection established")

    new_users_ch = make(chan *User, 10)
    new_topics_ch = make(chan *Topic, 10)
    new_messages_ch = make(chan *Message, 10)
    // go ExpAnswerBot()
    go GreeterBot()
    go RedditAnswerBot()

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

func GetHearthBeats() map[string]bool {
    ret := make(map[string]bool)

    key_name := redis_key_name("ExpAnswerBot")
    val, _ := redis_client.Get(key_name).Result()

    valint, err := strconv.ParseInt(val, 10, 64)

    if err == nil {
        hearthbeat_interval := time.Now().Unix() - valint
        dead_interval := int64(BotHearthBeatDead * time.Second)
        fmt.Println(hearthbeat_interval)
        fmt.Println(dead_interval)
        ret[key_name] = hearthbeat_interval < dead_interval
    } else {
        ret[key_name] = false
    }

    ret["lol"] = true

    return ret
}

func GreeterBot() {
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
        case new_user := <- new_users_ch:
            title := fmt.Sprintf("Welcome %s!", new_user.Username)
            user_detail_page := fmt.Sprintf("%s/users/%s", MakeBaseUrl(), new_user.Username)
            message := fmt.Sprintf(greeter_template, new_user.Username, user_detail_page)
            NewTopic(user, title, message)
        case <-time.After(10 * time.Second):
            func(){}()
        }
    }
}

func RedditAnswerBot() {
    bot_name := "RedditAnswerBot"
    user, _ := GetUser(bot_name)
    re := regexp.MustCompile(`\/r\/[a-zA-Z0-9_]+`)

    for {
        beat_hearth(bot_name)

        select {
        case new_message := <- new_messages_ch:
            if new_message.Author.Id != user.Id {
                matches := re.FindAllString(new_message.Message, -1)

                if len(matches) == 0{
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

                if new_message_text != "" {
                    NewMessage(user, new_message.Topic, new_message_text)
                }
            }
        case <-time.After(10 * time.Second):
            func(){}()
        }
    }
}

func ExpAnswerBot() {
    user, _ := GetUser("RiftBot")
    for {
        beat_hearth("ExpAnswerBot")

        select {
        case new_topic := <- new_topics_ch:
            NewMessage(user, new_topic, "Answer to a topic")
        case new_message := <- new_messages_ch:
            if new_message.Author.Id != user.Id {
                NewMessage(user, new_message.Topic, "Answer to a message")
            }
        case <-time.After(10 * time.Second):
            func(){}()
        }
    }
}

func redis_key_name(bot_name string) string {
    return fmt.Sprintf("%s_hearthbeat", bot_name)
}

func beat_hearth(bot_name string) {
    key_name := redis_key_name(bot_name)
    redis_client.Set(key_name, time.Now().Unix(), BotHearthBeatExpire * time.Second)
}
