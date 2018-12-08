package main

import (
    "time"

    "golang.org/x/crypto/bcrypt"
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

    Username string `sql:",unique"`
    PasswordHash string
    Signature string
    About string
    Usertype UserTypes

    Topics []*Topic
    Messages []*Message
}

func GenerateHash(password string) (string, error) {
    salted_bytes := []byte(password)
    hashed_bytes, err := bcrypt.GenerateFromPassword(salted_bytes, bcrypt.DefaultCost)

    if err != nil {
        return "", err
    }

    hash := string(hashed_bytes[:])
    return hash, nil
}

func NewUser(username string, user_type int, password string) *User {
    hash, err := GenerateHash(password)

    if err != nil {
        panic(err)
    }

    user := &User{
        Username: username,
        PasswordHash: hash,
        Usertype: Administrator,
    }

    return user
}

func (u User) GetUserType() string {
    if u.Usertype == Administrator {
        return "Administrator"
    } else if u.Usertype == Moderator {
        return "Moderator"
    } else if u.Usertype == Basic {
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

type InviteStatus uint8

const (
    Unused = iota
    Used
    Canceled
)

type Invite struct {
    DBObject
    Key string `sql:",unique"`
    Status InviteStatus
}

func (i Invite) GetInviteStatus() string {
    if i.Status == Unused {
        return "Unnused"
    } else if i.Status == Used {
        return "Used"
    } else if i.Status == Canceled {
        return "Canceled"
    }

    return "NoStatus"
}

func (i Invite) GetKeyUrl() string {
    return "TODO Key Url"
}
