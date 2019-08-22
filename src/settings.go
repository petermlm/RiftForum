package main

const (
    // Debug
    DebugMode = true

    // Server settings
    HostAndPort = ":8080"
    BaseUrl = "localhost:8080"
    Https = false

    // Database settings
    DatabaseConnRetries = 5

    DatabaseAddr = "postgres:5432"
    DatabaseDatabase = "riftforum_db"
    DatabaseUser = "riftforum_user"
    DatabasePassword = "riftforum_pass"

    // Redis settings
    RedisAddr = "redis:6379"

    // Bot settings
    BotHearthBeatPeriod = 10 // Seconds
    BotHearthBeatExpire = 60 // Seconds
    BotHearthBeatDead = 60 // Seconds
    BotChannelLag = 1024

    // Templating settings
    TemplatesDir = "templates"

    // Misc
    InviteSize = 12
    PageDefaultLimit = 20
    PageDefaultOffset = 0
)
