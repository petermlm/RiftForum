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

func SaveUser(user *User) {
    db := GetDBCon()
    err := db.Insert(user)

    if err != nil {
        panic(err)
    }
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

func InviteExists(invite_key string) bool {
    invite := new(Invite)
    err := db.Model(invite).
        Where("\"invite\".key = ?", invite_key).
        Select()

    if err != nil {
        return false
    }

    return true
}

func InviteSet(invite_key string, new_status InviteStatus) {
    invite := new(Invite)
    err := db.Model(invite).
        Where("\"invite\".key = ?", invite_key).
        Select()

    if err != nil {
        return
    }

    invite.Status = new_status

    db.Update(invite)
}

func InviteCancelAll() {
    db := GetDBCon()

    // TODO: User ORM instead of direct SQL query
    db.Model((*Invite)(nil)).Exec(`
        UPDATE invites
        SET status = 2
        WHERE status = 0
    `)
}

func GetTopics() []*Topic {
    db := GetDBCon()
    var topics []*Topic

    err := db.Model(&topics).
        Relation("Author").
        Relation("Messages", func(q *orm.Query) (*orm.Query, error) {
            return q.Order("message.created_at ASC"), nil
        }).
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
