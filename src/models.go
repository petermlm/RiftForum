package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type DBObject struct {
	Id        uint
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time `sql:"default:now()"`
}

/* ============================================================================
 * User
 * ============================================================================
 */

type UserTypes uint8

const (
	Administrator = iota + 1
	Moderator
	Basic
	Bot
)

type User struct {
	DBObject

	Username     string    `sql:",notnull,unique"`
	PasswordHash string    `sql:",notnull"`
	Signature    string    `sql:"default:''`
	About        string    `sql:"default:''`
	UserType     UserTypes `sql:",notnull"`
	Banned       bool      `sql:",defautl:false"`
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

func NewUser(username string, user_type UserTypes, password string) (*User, error) {
	hash, err := GenerateHash(password)

	if err != nil {
		panic(err)
	}

	if len(username) > Config.MaxUsernameSize {
		return nil, errors.New("Username is to big")
	}

	var valid_username = regexp.MustCompile("^[a-zA-Z0-9]{1,20}$")
	if !valid_username.MatchString(username) {
		return nil, errors.New("Invalid username")
	}

	user := &User{
		Username:     username,
		PasswordHash: hash,
		UserType:     user_type,
	}
	db.Insert(user)
	SendNewUser(user)
	return user, nil
}

func (u User) GetUserType() string {
	if u.UserType == Administrator {
		return "Administrator"
	} else if u.UserType == Moderator {
		return "Moderator"
	} else if u.UserType == Basic {
		return "Basic"
	} else if u.UserType == Bot {
		return "Bot"
	}
	return "NoType"
}

/* ============================================================================
 * Topic
 * ============================================================================
 */

const TitleLengthMax = 60

type Topic struct {
	DBObject

	Title string `sql:",notnull"`

	AuthorId uint  `sql:",notnull"`
	Author   *User `sql:",notnull"`

	Messages []*Message
}

func NewTopic(user *User, title_text string, message_text string) (*Topic, error) {
	var err error

	if len(title_text) > TitleLengthMax {
		return nil, errors.New("Title is too big")
	}

	// Topic
	topic := &Topic{
		Title:    title_text,
		Author:   user,
		AuthorId: user.Id,
	}

	err = db.Insert(topic)

	if err != nil {
		panic(err)
	}

	// Message
	message := &Message{
		Message:  message_text,
		Author:   user,
		AuthorId: user.Id,
		Topic:    topic,
		TopicId:  topic.Id,
	}

	err = db.Insert(message)

	if err != nil {
		panic(err)
	}

	SendNewMessage(message)
	SendNewTopic(topic)
	return topic, nil
}

/* ============================================================================
 * Message
 * ============================================================================
 */

type Message struct {
	DBObject
	Message string `sql:",notnull"`

	AuthorId uint  `sql:",notnull"`
	Author   *User `sql:",notnull"`

	TopicId uint   `sql:",notnull"`
	Topic   *Topic `sql:",notnull"`
}

func NewMessage(user *User, topic *Topic, message_text string) *Message {
	UpdateTopic(topic)

	message := &Message{
		Message:  message_text,
		Author:   user,
		AuthorId: user.Id,
		Topic:    topic,
		TopicId:  topic.Id,
	}

	err := db.Insert(message)

	if err != nil {
		panic(err)
	}

	SendNewMessage(message)
	return message
}

/* ============================================================================
 * Invite
 * ============================================================================
 */

type InviteStatus uint8

const (
	Unused = iota + 1
	Used
	Canceled
)

type Invite struct {
	DBObject
	Key    string       `sql:",notnull,unique"`
	Status InviteStatus `sql:",notnull"`
}

func (i Invite) GetInviteStatus() string {
	if i.Status == Unused {
		return "Unused"
	} else if i.Status == Used {
		return "Used"
	} else if i.Status == Canceled {
		return "Canceled"
	}

	return "NoStatus"
}

func (i Invite) GetKeyURL() string {
	base_url := MakeBaseURL()
	return fmt.Sprintf("%s/register?key=%s", base_url, i.Key)
}
