package main

const (
    // Debug Mode
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

    DatabaseTemp = true

    // Templating settings
    TemplatesDir = "templates"

    // Misc
    InviteSize = 12
    PageDefaultLimit = 20
    PageDefaultOffset = 0
)
