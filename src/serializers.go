package main

import (
    "log"
    "strings"
    "html/template"

    "github.com/frustra/bbcode"
)

var bbcode_compiler bbcode.Compiler

func InitSers() {
    // Booleans: autoCloseTags, ignoreUnmatchedClosingTags
    bbcode_compiler = bbcode.NewCompiler(true, true)
    log.Println("Serializers initialized")
}

type UserInfo struct {
    Username string
    Usertype UserTypes
}

func (u *UserInfo) IsAdmin() bool {
    return u.Usertype == Administrator
}

func (u *UserInfo) IsMod() bool {
    return u.Usertype == Administrator || u.Usertype == Moderator
}

type RiftDataI interface {
    SetUserInfo(UserInfo *UserInfo)
    SetPath(path string)
    SetBannerSentence(sentence string)
    HasUser() bool
    IsAdmin() bool
}

type RiftData struct {
    UserInfo *UserInfo
    Path string
    BannerSentence string
}

func (r *RiftData) SetUserInfo(UserInfo *UserInfo) {
    r.UserInfo = UserInfo
}

func (r *RiftData) SetPath(path string) {
    r.Path = path
}

func (r *RiftData) SetBannerSentence(sentence string) {
    r.BannerSentence = sentence
}

func (r *RiftData) HasUser() bool {
    return r.UserInfo != nil
}

func (r *RiftData) IsAdmin() bool {
    if r.UserInfo == nil {
        return false
    }

    return r.UserInfo.IsAdmin()
}

func (r *RiftData) IsMod() bool {
    if r.UserInfo == nil {
        return false
    }

    return r.UserInfo.IsMod()
}

type EmptyData struct {
    RiftData
}

type LoginData struct {
    RiftData
    NextPage string
}

type UserListData struct {
    Username string
    Usertype string
    CreatedAt string
}

type UsersListData struct {
    RiftData
    PageNum int
    PageSize int
    PageMax int
    Users []UserListData
}

type UserData struct {
    RiftData
    Username string
    Usertype string
    CreatedAt string
    AboutParagraphs []template.HTML
    About string
    SignatureParagraphs []template.HTML
    Signature string
    Banned bool
    SameUser bool
}

func (u *UserData) IsBanned() bool {
    return u.Banned
}

func (u *UserData) IsSameUser() bool {
    return u.SameUser
}

type RegisterData struct {
    RiftData
    Key string
}

func (r *RegisterData) HasKey() bool {
    return r.Key != ""
}

type ChangePasswordData struct {
    RiftData
    Username string
    is_for_admin bool
    render_error bool
}

func (c *ChangePasswordData) IsForAdmin() bool {
    return c.is_for_admin
}

func (c *ChangePasswordData) RenderError() bool {
    return c.render_error
}

type InviteListData struct {
    Key string
    Status string
    CreatedAt string
}

type InvitesListData struct {
    RiftData
    PageNum int
    PageSize int
    PageMax int
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
    PageNum int
    PageSize int
    PageMax int
    Topics []TopicListData
}

type MessageData struct {
    AuthorId uint
    AuthorUsername string
    AuthorUsertype string
    SignatureParagraphs []template.HTML
    CreatedAt string
    MessageParagraphs []template.HTML
}

type TopicData struct {
    RiftData
    TopicId uint
    Title string
    Messages []*MessageData
    PageNum int
    PageSize int
    PageMax int
}

func SerializeEmpty() *EmptyData {
    return new(EmptyData)
}

func SerializeLogin(next_page string) *LoginData {
    ser_login := new(LoginData)

    if next_page == "" {
        ser_login.NextPage = "/"
    } else {
        ser_login.NextPage = next_page
    }

    return ser_login
}

func SerializeUsers(users []*User, page Page) *UsersListData {
    ser_users := new(UsersListData)

    ser_users.PageNum = page.get_num()
    ser_users.PageSize = page.get_size()
    ser_users.PageMax = CountUsersPages(page)

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

func SerializeUser(cusername string, user *User) *UserData {
    ser_user := new(UserData)

    ser_user.Username = user.Username
    ser_user.Usertype = user.GetUserType()
    ser_user.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
    ser_user.AboutParagraphs = stringToOutputHtml(user.About)
    ser_user.About = user.About
    ser_user.SignatureParagraphs = stringToOutputHtml(user.Signature)
    ser_user.Signature = user.Signature
    ser_user.Banned = user.Banned
    ser_user.SameUser = cusername == user.Username

    return ser_user
}

func SerializeRegister(key string) *RegisterData {
    ser_register := new(RegisterData)
    ser_register.Key = key
    return ser_register
}

func SerializeChangePassword(user *User, is_for_admin bool, render_error bool) *ChangePasswordData {
    ser_change_password := new(ChangePasswordData)
    ser_change_password.Username = user.Username
    ser_change_password.is_for_admin = is_for_admin
    ser_change_password.render_error = render_error
    return ser_change_password
}

func SerializeInvites(invites []*Invite, page Page) *InvitesListData {
    ser_invites := new(InvitesListData)

    ser_invites.PageNum = page.get_num()
    ser_invites.PageSize = page.get_size()
    ser_invites.PageMax = CountInvitesPages(page)

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

func SerializeTopics(topics []*Topic, page Page) *TopicsListData {
    ser_topics := new(TopicsListData)

    ser_topics.PageNum = page.get_num()
    ser_topics.PageSize = page.get_size()
    ser_topics.PageMax = CountTopicsPages(page)

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

func SerializeTopic(topic *Topic, page Page) *TopicData {
    var messages []*MessageData

    for _, message := range topic.Messages {
        message_struct := &MessageData {
            AuthorId: message.Author.Id,
            AuthorUsername: message.Author.Username,
            AuthorUsertype: message.Author.GetUserType(),
            SignatureParagraphs: stringToOutputHtml(message.Author.Signature),
            CreatedAt: message.CreatedAt.Format("2006-01-02 15:04:05"),
            MessageParagraphs: stringToOutputHtml(message.Message),
        }

        messages = append(messages, message_struct)
    }

    ser_topic := &TopicData {
        TopicId: topic.Id,
        Title: topic.Title,
        Messages: messages,
        PageNum: page.get_num(),
        PageSize: page.get_size(),
        PageMax: CountMessagePages(topic.Id, page),
    }

    return ser_topic
}

func stringToOutputHtml(str string) []template.HTML {
    if str == "" {
        return make([]template.HTML, 0)
    }

    str_pars := strings.Split(str, "\r\n")
    str_pars_html := make([]template.HTML, len(str_pars))

    for i, str_par := range str_pars {
        str_pars_html[i] = template.HTML(bbcode_compiler.Compile(str_par))
    }

    return str_pars_html
}
