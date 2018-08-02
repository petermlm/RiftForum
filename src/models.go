package main

import (
    "time"
)

type DBObject struct {
    Id uint
    Created time.Time
    Modified time.Time
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

    AuthorID uint
    Author *User
}

/* ============================================================================
 * Message
 * ============================================================================
 */

type Message struct {
    DBObject
    Message string

    AuthorID uint
    Author *User

    TopicID uint
    Topic *User
}

/* ============================================================================
 * Invite
 * ============================================================================
 */

type Invite struct {
    DBObject
    Key string
    Userd bool
}
