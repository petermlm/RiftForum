package main

const (
    // Debug
    DebugMode = true

    // Server
    HostAndPort = ":8080"
    BaseUrl = "localhost:8080"
    Https = false

    // Database and Reddis
    DatabaseConnRetries = 5
    DatabaseAddr = "postgres:5432"
    DatabaseDatabase = "riftforum_db"
    DatabaseUser = "riftforum_user"
    DatabasePassword = "riftforum_pass"
    RedisAddr = "redis:6379"

    // Users and Invites
    AdminUsername = "admin"
    FirstUsername = "Pedro"
    DefaultPassword = "pl"
    InviteSize = 12

    // Bots
    BotHearthBeatPeriod = 10 // Seconds
    BotHearthBeatExpire = 60 // Seconds
    BotHearthBeatDead = 60 // Seconds
    BotChannelLag = 1024

    // Templating
    TemplatesDir = "templates"

    // Pages
    PageDefaultLimit = 20
    PageDefaultOffset = 0
    PageDefaultSize = 20
    PageDefaultNum = 0
)
