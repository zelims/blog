package app

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-github/github"
	"github.com/jmoiron/sqlx"
	"github.com/revel/revel"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

var DB *sqlx.DB

func initDB() {
	driver := revel.Config.StringDefault("db.driver", "mysql")
	connectString := revel.Config.StringDefault("db.connect", "root:@(localhost:3306)/blog")

	db, err := sqlx.Connect(driver, connectString)
	if err != nil {
		log.Fatal("[!] DB Err: ", err)
	}

	DB = db
}

func setupTemplateFuncs() {
	revel.TemplateFuncs["strcat"] = func(strs ...string) string {
		return strings.Trim(strings.Join(strs, ""), " ")
	}
	revel.TemplateFuncs["strcmp"] = func(str1, str2 string) bool {
		return str1 == str2
	}
	revel.TemplateFuncs["strintcmp"] = func(str string, i int) bool {
		itoa := strconv.Itoa(i)
		return str == itoa
	}
	revel.TemplateFuncs["arrsize"] = func(strs []string) int {
		return len(strs)
	}
	revel.TemplateFuncs["printtags"] = func(tags []string, all bool) template.HTML {
		output := "<small><i class=\"fas fa-tags\"></i>  "
		for i, tag := range tags {
			if i < len(tags) && i != 0 {
				output += ", "
			}
			output += fmt.Sprintf("<a href=\"/tag/%s\">%s</a>", tag, tag)
			if i == 2 && !all {
				break
			}
		}
		output += "</small>"
		return template.HTML(output)
	}
	revel.TemplateFuncs["md"] = func(str string) template.HTML {
		return template.HTML(string(blackfriday.MarkdownCommon([]byte(str))))
	}
	revel.TemplateFuncs["html"] = func(str string) template.HTML {
		return template.HTML(str)
	}
	revel.TemplateFuncs["bannerImg"] = func(id int, imgName string) template.HTML {
		img := fmt.Sprintf("/public/images/posts/%s.jpg", imgName)
		if _, err := os.Stat(revel.BasePath + img); err != nil {
			return ""
		}
		return template.HTML("<style>header.page-header::before { background: no-repeat center url(" + img + ");}</style>")
	}
	revel.TemplateFuncs["printImages"] = func(images []string) template.HTML {
		out := "printing the beautiful images"
		// loop through images
			// output img tag

		return template.HTML(out)
	}
	revel.TemplateFuncs["pagination"] = func(n int)  (stream chan int){
		stream = make(chan int)
		go func() {
			for i := 1; i <= n; i++ {
				stream <- i
			}
			close(stream)
		}()
		return
	}
	revel.TemplateFuncs["github_time_format"] = func(t *github.Timestamp) string {
		return t.Format("02 Jan 2006 15:04")
	}
	revel.TemplateFuncs["time_format"] = func(unixTime int) string {
		return time.Unix(int64(unixTime), 0).Format("02 Jan 2006 15:04:05")
	}
	revel.TemplateFuncs["itoa"] = func(i int) string {
		return strconv.Itoa(i)
	}
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)

	revel.OnAppStart(initDB)
	revel.OnAppStart(setupTemplateFuncs)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
