package main

const (
    // Server settings
    HostAndPort = ":8080"

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
    PageDefaultLimit = 20
    PageDefaultOffset = 0
)
