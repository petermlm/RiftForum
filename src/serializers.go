package main

func SerializeTopics(topics []Topic) []interface{} {
    var ser_topics []interface{}

    for _, topic := range topics {
        msg_count := len(topic.Messages)
        last_message := topic.Messages[msg_count - 1]

        ser_topic := struct {
            Title string
            AuthorId uint
            AuthorUsername string
            CreatedAt string
            MessageCount int
            LastAuthor string
            LastTimestamp string
        }{
            Title: topic.Title,
            AuthorId: topic.Author.Id,
            AuthorUsername: topic.Author.Username,
            CreatedAt: topic.CreatedAt.Format("2006-01-02 15:04:05"),
            MessageCount: msg_count,
            LastAuthor: last_message.Author.Username,
            LastTimestamp: last_message.CreatedAt.Format("2006-01-02 15:04:05"),
        }

        ser_topics = append(ser_topics, ser_topic)
    }

    return ser_topics
}
