package main

import (
    "fmt"
    "time"
    "strconv"

    redis "github.com/go-redis/redis"
)

var redis_client *redis.Client
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

    new_topics_ch = make(chan *Topic, 10)
    new_messages_ch = make(chan *Message, 10)
    go ExpAnswerBot()
    fmt.Println("Bots Started")
}

func SendNewTopic(topic *Topic) {
    new_topics_ch <- topic
}

func SendNewMessage(message *Message) {
    new_messages_ch <- message
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

func ExpAnswerBot() {
    user := get_rift_bot_user()
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

func get_rift_bot_user() *User {
    user := new(User)

    err := db.Model(user).Where("Username = 'RiftBot'").Select()

    if err != nil {
        panic(err)
    }

    return user
}

func redis_key_name(bot_name string) string {
    return fmt.Sprintf("%s_hearthbeat", bot_name)
}

func beat_hearth(bot_name string) {
    key_name := redis_key_name(bot_name)
    redis_client.Set(key_name, time.Now().Unix(), BotHearthBeatExpire * time.Second)
}
