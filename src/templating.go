package main

import (
    "strings"
    "log"
    "os"
    "io/ioutil"
    "net/http"
    "html/template"
)

var templates *template.Template

func tpl_dir(tpl_name string) string {
    return "../templates/" + tpl_name
}

func InitTmpl() {
    var all_files []string
    var files []os.FileInfo
    var err error

    files, err = ioutil.ReadDir("../templates")
    if err != nil {
        log.Println("can't read directory with templates")
        log.Fatal(err)
    }

    for _, file := range files {
        filename := file.Name()
        if strings.HasSuffix(filename, ".html") {
            all_files = append(all_files, "../templates/" + filename)
        }
    }

    templates, err = template.New("webpage").ParseFiles(all_files...)

    if err != nil {
        log.Println("Can't parse template files")
        log.Fatal(err)
    }
}

func Render(writer *http.ResponseWriter, r *http.Request, tpl_name string, data RiftDataI) {
    ctx := r.Context()
    user_info, ok := ctx.Value("UserInfo").(*UserInfo)

    if ok {
        data.SetUserInfo(user_info)
    }

    err := templates.ExecuteTemplate(*writer, tpl_name, data)

    if err != nil {
        log.Println("Can't execute template")
        log.Fatal(err)
    }
}
