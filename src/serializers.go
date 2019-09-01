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
    AddCustomBBCode(&bbcode_compiler)
    log.Println("Serializers initialized")
}

type UserInfo struct {
    Id uint
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
    SetData(path string)
    HasUser() bool
    IsAdmin() bool
}

type RiftData struct {
    UserInfo *UserInfo
    Path string
    BannerSentence string
    ApiBase string
    VersionString string
}

func (r *RiftData) SetUserInfo(UserInfo *UserInfo) {
    r.UserInfo = UserInfo
}

func (r *RiftData) SetData(path string) {
    r.Path = path
    r.BannerSentence = get_banner_sentence()
    r.ApiBase = ApiBase
    r.VersionString = VersionString
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
    Banned bool
}

func (u *UserListData) IsBanned() bool {
    return u.Banned
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

type RegisterErrors struct {
    invite_key_bad bool
    username_alreay_taken bool
    username_is_invalid bool
    passwords_dont_match bool
}

type RegisterData struct {
    RiftData
    Key string
    errors RegisterErrors
}

func (r *RegisterData) HasKey() bool {
    return r.Key != ""
}

func (r *RegisterData) InviteKeyBad() bool {
    return r.errors.invite_key_bad
}

func (r *RegisterData) UsernameAlreayTaken() bool {
    return r.errors.username_alreay_taken
}

func (r *RegisterData) UsernameIsInvalid() bool {
    return r.errors.username_is_invalid
}

func (r *RegisterData) PasswordsDontMatch() bool {
    return r.errors.passwords_dont_match
}

type ChangePasswordErrors struct {
    old_password_wrong bool
    new_passwords_not_equal bool
}

type ChangePasswordData struct {
    RiftData
    Username string
    is_for_admin bool
    errors ChangePasswordErrors
}

func (c *ChangePasswordData) IsForAdmin() bool {
    return c.is_for_admin
}

func (c *ChangePasswordData) OldPasswordIsWrong() bool {
    return c.errors.old_password_wrong
}

func (c *ChangePasswordData) NewPasswordsNotEqual() bool {
    return c.errors.new_passwords_not_equal
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

type BotsData struct {
    RiftData
    HearthBeatStatus map[string]bool
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
    Id uint
    AuthorId uint
    AuthorUsername string
    AuthorUsertype string
    AuthorBanned bool
    SignatureParagraphs []template.HTML
    CreatedAt string
    MessageParagraphs []template.HTML
}

func (m *MessageData) IsAuthorBanned() bool {
    return m.AuthorBanned
}

func (m *MessageData) CanEdit() bool {
    return true
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

type MessageEditData struct {
    RiftData
    Id uint
    TopicId uint
    Title string
    Message string
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
            Banned: user.Banned,
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

func SerializeRegister(key string, errors RegisterErrors) *RegisterData {
    ser_register := new(RegisterData)
    ser_register.Key = key
    ser_register.errors = errors
    return ser_register
}

func SerializeChangePassword(
    user *User,
    is_for_admin bool,
    errors ChangePasswordErrors,
) *ChangePasswordData {
    ser_change_password := new(ChangePasswordData)
    ser_change_password.Username = user.Username
    ser_change_password.is_for_admin = is_for_admin
    ser_change_password.errors = errors
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

func SerializeBots(hearthbeat_status map[string]bool) *BotsData {
    ser_bots := new(BotsData)
    ser_bots.HearthBeatStatus = hearthbeat_status
    return ser_bots
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
            Id: message.Id,
            AuthorId: message.Author.Id,
            AuthorUsername: message.Author.Username,
            AuthorUsertype: message.Author.GetUserType(),
            AuthorBanned: message.Author.Banned,
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

func SerializeMessageEdit(message *Message) *MessageEditData {
    data := new(MessageEditData)
    data.Id = message.Id
    data.TopicId = message.Topic.Id
    data.Title = message.Topic.Title
    data.Message = message.Message
    return data
}

func stringToOutputHtml(str string) []template.HTML {
    if str == "" {
        return make([]template.HTML, 0)
    }

    str_bbcoded := bbcode_compiler.Compile(str)
    str_pars := strings.Split(str_bbcoded, "\r\n")
    str_pars_html := make([]template.HTML, len(str_pars))

    for i, str_par := range str_pars {
        str_pars_html[i] = template.HTML(str_par)
    }

    return str_pars_html
}
