package main

func GetTopics() []Topic {
    db := GetDBCon()
    var topics []Topic

    err := db.Model(&topics).
        Relation("Author").
        Relation("Messages").
        Relation("Messages.Author").
        Limit(50).
        Select()

    if err != nil {
        panic(err)
    }

    return topics
}
