package main

import (
    "time"
)

type DBObject struct {
    Id uint
    CreatedAt time.Time `sql:"default:now()"`
    UpdatedAt time.Time `sql:"default:now()"`
}

/* ============================================================================
 * User
 * ============================================================================
 */

type UserTypes uint8

const (
    Administrator = iota
    Moderator
    Basic
)

type User struct {
    DBObject

    Username string
    PasswordHash string
    Signature string
    About string
    UserType UserTypes

    Topics []*Topic
    Messages []*Message
}

func NewUser(username string, user_type string, password string) User {
    user := User{
        Username: "admin",
        UserType: Administrator,
    }

    return user
}

func (u User) GetUserType() string {
    if u.UserType == Administrator {
        return "Administrator"
    } else if u.UserType == Moderator {
        return "Moderator"
    } else if u.UserType == Basic {
        return "Basic"
    }

    return "NoType"
}

/* ============================================================================
 * Topic
 * ============================================================================
 */

type Topic struct {
    DBObject

    Title string

    AuthorId uint
    Author *User

    Messages []*Message
}

/* ============================================================================
 * Message
 * ============================================================================
 */

type Message struct {
    DBObject
    Message string

    AuthorId uint
    Author *User

    TopicId uint
    Topic *Topic
}

/* ============================================================================
 * Invite
 * ============================================================================
 */

type Invite struct {
    DBObject
    Key string
    Used bool
}
