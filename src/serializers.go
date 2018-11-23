package main

type RiftDataI interface {
    SetUserInfo(UserInfo *UserInfo)
    HasUser() bool
}

type RiftData struct {
    UserInfo *UserInfo
}

func (r *RiftData) SetUserInfo(UserInfo *UserInfo) {
    r.UserInfo = UserInfo
}

func (r *RiftData) HasUser() bool {
    return r.UserInfo != nil
}

type TopicListData struct {
    TopicId uint
    Title string
    AuthorId uint
    AuthorUsername string
    CreatedAt string
    MessageCount int
    LastAuthor string
    LastTimestamp string
}

type TopicsListData struct {
    RiftData
    Topics []TopicListData
}

type MessageData struct {
    AuthorId uint
    AuthorUsername string
    CreatedAt string
    Message string
}

type TopicData struct {
    RiftData
    TopicId uint
    Title string
    Messages []*MessageData
}

func SerializeTopics(topics []*Topic) *TopicsListData {
    ser_topics := new(TopicsListData)

    for _, topic := range topics {
        msg_count := len(topic.Messages)
        last_message := topic.Messages[msg_count - 1]

        ser_topic := TopicListData {
            TopicId: topic.Id,
            Title: topic.Title,
            AuthorId: topic.Author.Id,
            AuthorUsername: topic.Author.Username,
            CreatedAt: topic.CreatedAt.Format("2006-01-02 15:04:05"),
            MessageCount: msg_count,
            LastAuthor: last_message.Author.Username,
            LastTimestamp: last_message.CreatedAt.Format("2006-01-02 15:04:05"),
        }

        ser_topics.Topics = append(ser_topics.Topics, ser_topic)
    }

    return ser_topics
}

func SerializeTopic(topic *Topic) *TopicData {
    var messages []*MessageData

    for _, message := range topic.Messages {
        message_struct := &MessageData {
            AuthorId: message.Author.Id,
            AuthorUsername: message.Author.Username,
            CreatedAt: message.CreatedAt.Format("2006-01-02 15:04:05"),
            Message: message.Message,
        }

        messages = append(messages, message_struct)
    }

    ser_topic := &TopicData {
        TopicId: topic.Id,
        Title: topic.Title,
        Messages: messages,
    }

    return ser_topic
}
