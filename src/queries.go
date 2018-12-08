package main

import (
    "errors"
    "time"
    "github.com/go-pg/pg/orm"
)

func GetUser(username string) (*User, error) {
    db := GetDBCon()

    user := new(User)
    err := db.Model(user).
        Where("\"user\".username = ?", username).
        Select()

    if err != nil {
        return nil, errors.New("User doesn't exist")
    }

    return user, nil
}

func GetInvites() []*Invite {
    db := GetDBCon()
    var invites []*Invite

    err := db.Model(&invites).
        Order("invite.created_at DESC").
        Limit(10).
        Select()

    if err != nil {
        panic(err)
    }

    return invites
}

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