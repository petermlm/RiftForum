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
