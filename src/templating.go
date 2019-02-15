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

    funcMap := template.FuncMap{
        "add": func(a, b int) int {
            return a + b
        },
        "sub": func(a, b int) int {
            return a - b
        },
    }

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

    templates, err = template.
        New("webpage").
        Funcs(funcMap).
        ParseFiles(all_files...)

    if err != nil {
        log.Println("Can't parse template files")
        log.Fatal(err)
    }
}

func Render(res *http.ResponseWriter, req *http.Request, tpl_name string, data RiftDataI) {
    ctx := req.Context()
    user_info, ok := ctx.Value("UserInfo").(*UserInfo)

    if ok {
        data.SetUserInfo(user_info)
    }

    data.SetPath(req.URL.Path)

    err := templates.ExecuteTemplate(*res, tpl_name, data)

    if err != nil {
        log.Println("Can't execute template")
        log.Fatal(err)
    }
}

func Redirect(res *http.ResponseWriter, req *http.Request, url string) {
    http.Redirect(*res, req, url, http.StatusSeeOther)
}

func Login(res *http.ResponseWriter, req *http.Request) {
    redirect := "/login"
    Redirect(res, req, redirect)
}

func NotFound(res *http.ResponseWriter, req *http.Request) {
    Render(res, req, "404.html", SerializeEmpty())
}

func AdminOnly(res *http.ResponseWriter, req *http.Request) {
    Render(res, req, "admin_only.html", SerializeEmpty())
}

func OperationNotAllowed(res *http.ResponseWriter, req *http.Request) {
    Render(res, req, "operation_not_allowed.html", SerializeEmpty())
}
