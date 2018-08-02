package main

import (
    "log"
    "net/http"
    "html/template"
)

func tpl_dir(tpl_name string) string {
    return "../templates/" + tpl_name
}

func Render(tpl_name string, w http.ResponseWriter, data interface{}) {
    t, err := template.New("webpage").ParseFiles(tpl_dir(tpl_name))

    if err != nil {
        log.Fatal("Can't open template")
    }

    err = t.ExecuteTemplate(w, tpl_name, data)

    if err != nil {
        log.Fatal(err)
    }
}
