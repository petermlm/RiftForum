package main

type UserInfo struct {
    Username string
    Usertype UserTypes
}

type RiftDataI interface {
    SetUserInfo(UserInfo *UserInfo)
    HasUser() bool
    IsAdmin() bool
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

func (r *RiftData) IsAdmin() bool {
    if r.UserInfo == nil {
        return false
    }

    return r.UserInfo.Usertype == Administrator
}

func (r *RiftData) IsMod() bool {
    if r.UserInfo == nil {
        return false
    }

    return r.UserInfo.Usertype == Administrator || r.UserInfo.Usertype == Moderator
}

type EmptyData struct {
    RiftData
}

type UserListData struct {
    Username string
    Usertype string
    CreatedAt string
}

type UsersListData struct {
    RiftData
    Users []UserListData
}

type UserData struct {
    RiftData
    Username string
    Usertype string
    CreatedAt string
    About string
    AboutF string
    Signature string
    SignatureF string
}

type RegisterData struct {
    RiftData
    Key string
}

func (r *RegisterData) HasKey() bool {
    return r.Key != ""
}

type InviteListData struct {
    Key string
    Status string
    CreatedAt string
}

type InvitesListData struct {
    RiftData
    Invites []InviteListData
}

type InviteNewData struct {
    RiftData
    Key string
    KeyUrl string
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
    AuthorUsertype string
    SignatureF string
    CreatedAt string
    Message string
}

type TopicData struct {
    RiftData
    TopicId uint
    Title string
    Messages []*MessageData
}

func SerializeEmpty() *EmptyData {
    return new(EmptyData)
}

func SerializeUsers(users []*User) *UsersListData {
    ser_users := new(UsersListData)

    for _, user := range users {
        ser_user := UserListData {
            Username: user.Username,
            Usertype: user.GetUserType(),
            CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
        }

        ser_users.Users = append(ser_users.Users, ser_user)
    }

    return ser_users
}

func SerializeUser(user *User) *UserData {
    ser_user := new(UserData)

    ser_user.Username = user.Username
    ser_user.Usertype = user.GetUserType()
    ser_user.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
    ser_user.AboutF = user.About
    ser_user.About = user.About
    ser_user.SignatureF = user.Signature
    ser_user.Signature = user.Signature

    return ser_user
}

func SerializeRegister(key string) *RegisterData {
    ser_register := new(RegisterData)
    ser_register.Key = key
    return ser_register
}

func SerializeInvites(invites []*Invite) *InvitesListData {
    ser_invites := new(InvitesListData)

    for _, invite := range invites {
        ser_invite := InviteListData {
            Key: invite.Key,
            Status: invite.GetInviteStatus(),
            CreatedAt: invite.CreatedAt.Format("2006-01-02 15:04:05"),
        }

        ser_invites.Invites = append(ser_invites.Invites, ser_invite)
    }

    return ser_invites
}

func SerializeInviteNew(new_invite *Invite) *InviteNewData {
    return &InviteNewData {
        Key: new_invite.Key,
        KeyUrl: new_invite.GetKeyUrl(),
    }
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
            AuthorUsertype: message.Author.GetUserType(),
            SignatureF: message.Author.Signature,
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
