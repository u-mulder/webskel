package main

import (
	"fmt"
	"os"
)

var tplFileContents = `{{define "IndexTpl"}}
<!doctype html>
<html>
    <head>
        <title>Title Here</title>
    </head>

    <body>
        <div>
            Content here
        </div>

        <script type="text/javascript" src="script.js"></script>
        <script type="text/javascript"> jsCode Here </script>
    </body>
</html>
{{end}}`

var goFileContents = `package main

import (
    "database/sql"
    "fmt"
    "github.com/julienschmidt/httprouter"
    _ "github.com/lib/pq"
    "html/template"
    "net/http"
    "os"
)

var db *sql.DB
var tplPath = "public/assets/tpls"
var templates *template.Template

func init() {
    var err error
    db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        panic(err)
    }

    templates = template.Must(template.ParseGlob(tplPath + "/*.tpl"))
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    r := httprouter.New()
    /* serving static files in /assets/ path */
    r.ServeFiles("/assets/*filepath", http.Dir("public/assets/"))

    r.GET("/", indexPage)
    r.GET("/otherUrl", otherPage)
    r.POST("/otherUrl", postOtherPage)

    fmt.Println("Listening server on port " + port)
    http.ListenAndServe(":"+port, r)

    defer db.Close()
}

func indexPage(rw http.ResponseWriter, rq *http.Request, _ httprouter.Params) {
    var tplData = struct{}{}
    err := templates.ExecuteTemplate(rw, "IndexTpl", tplData)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }
}

func otherPage(rw http.ResponseWriter, rq *http.Request, _ httprouter.Params) {
    // TODO
}

func postOtherPage(rw http.ResponseWriter, rq *http.Request, _ httprouter.Params) {
    // TODO
}`

func main() {
	// Finding current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error detecting current directory: ", err)
		os.Exit(1)
	}

	// Create 'main.go' file
	f, err := os.Create(cwd + "/main.go")
	if err != nil {
		fmt.Println("Error creating 'main.go' file: ", err)
	}
	defer f.Close()

	_, err = f.WriteString(goFileContents)
	if err != nil {
		fmt.Println("Error writing to 'main.go' file: ", err)
	}
	f.Sync()

	// Create public folders
	publicPath := cwd + "/public/assets/"
	subfolders := []string{"css", "images", "js", "tpls"}
	for _, v := range subfolders {
		err := os.MkdirAll(publicPath+v, 0755)
		if err != nil {
			fmt.Println("Error creating folder '", v, "': ", err)
		} else {
			if v == "tpls" {
				// Create 'index.tpl' file
				f, err := os.Create(publicPath + v + "/index.tpl")
				if err != nil {
					fmt.Println("Error creating 'index.tpl' file: ", err)
				} else {
					_, err = f.WriteString(tplFileContents)
					if err != nil {
						fmt.Println("Error writing to 'index.tpl' file: ", err)
					}
					f.Sync()
				}
				defer f.Close()
			}
		}
	}
}
