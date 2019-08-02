package main

import (
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path"
    "strings"
)

var templates *template.Template

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

    files, err = ioutil.ReadDir(TemplatesDir)
    if err != nil {
        RiftForumPanic("can't read directory with templates", err)
    }

    for _, file := range files {
        filename := file.Name()
        if strings.HasSuffix(filename, ".html") {
            all_files = append(all_files, tpl_dir(filename))
        }
    }

    templates, err = template.
        New("webpage").
        Funcs(funcMap).
        ParseFiles(all_files...)

    if err != nil {
        RiftForumPanic("Can't parse template files", err)
    }

    log.Println("Templating initialized")
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
        RiftForumPanic("Can't execute template", err)
    }
}

func Redirect(res *http.ResponseWriter, req *http.Request, url string) {
    http.Redirect(*res, req, url, http.StatusSeeOther)
}

func Login(res *http.ResponseWriter, req *http.Request) {
    Redirect(res, req, "/login")
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

func tpl_dir(tpl_name string) string {
    return path.Join(TemplatesDir, tpl_name)
}
