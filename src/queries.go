package main

import (
    "time"
    "github.com/go-pg/pg/orm"
)

func GetTopics() []*Topic {
    db := GetDBCon()
    var topics []*Topic

    err := db.Model(&topics).
        Relation("Author").
        Relation("Messages").
        Relation("Messages.Author").
        Order("topic.updated_at DESC").
        Limit(50).
        Select()

    if err != nil {
        panic(err)
    }

    return topics
}

func GetTopic(topic_id uint) *Topic {
    db := GetDBCon()

    topic := new(Topic)
    err := db.Model(topic).
        Relation("Author").
        Relation("Messages", func(q *orm.Query) (*orm.Query, error) {
            return q.Order("message.created_at ASC"), nil
        }).
        Relation("Messages.Author").
        Where("topic.id = ?", topic_id).
        Select()

    if err != nil {
        panic(err)
    }

    return topic
}

func UpdateTopic(topic *Topic) {
    topic.UpdatedAt = time.Now()
    err := db.Update(topic)

    if err != nil {
        panic(err)
    }
}
