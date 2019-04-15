package models

import (
	"fmt"
	"github.com/avct/uasurfer"
	"github.com/oschwald/geoip2-golang"
	"github.com/revel/revel"
	"github.com/satori/go.uuid"
	"github.com/zelims/blog/app/database"
	"log"
	"net"
	"strings"
	"time"
)

type CountryAnalytics struct {
	CountryCode		string	`db:"country"`
	UserCount		int		`db:"count"`
}

type AnalyticData struct {
	ID			int		`db:"ID"`
	UUID		string	`db:"uuid"`
	Page		string
	IP			string	`db:"ip_address"`
	Country		string
	OS			string
	Browser		string
	Device		string
	Time		int
	Visits		int		`db:"visits"`
}

type AnalyticsPosts struct {
	Post
	Count		int		`db:"count"`
}

func TrackUser(c *revel.Controller) revel.Result {
	trackID := c.Session["trackID"]
	if trackID == nil {
		u, err := uuid.NewV4()
		if err != nil {
			log.Printf("UUID Error: %s", err.Error())
			return c.RenderTemplate("errors/500.html")
		}
		c.Session["trackID"] = u.String()
	}
	request := c.Request
	ipAddress := strings.Split(request.Header.Get("X-Forwarded-For"), ":")[0]
	if ipAddress == "127.0.0.1" || ipAddress == "localhost" {
		//return nil
	}

	countryCode := getISOCodeFromIP(ipAddress)
	if countryCode == "" {
		countryCode = "Unknown"
	}

	ua := uasurfer.Parse(request.UserAgent())

	OS := strings.Split(ua.OS.Name.String(), "OS")[1]
	browser := strings.Title(strings.Split(strings.ToLower(ua.Browser.Name.String()), "browser")[1])
	device := strings.Title(strings.Split(strings.ToLower(ua.DeviceType.String()), "device")[1])

	_, err := database.Handle.NamedExec(`INSERT INTO analytics(ID,uuid,page,ip_address,country,os,browser,device,time)` +
		` VALUES(:id,:uuid,:page,:ip,:country,:os,:browser,:device,:time)`,
		map[string]interface{}{
			"id":		nil,
			"uuid":		trackID,
			"page":		request.URL.String(),
			"ip":		ipAddress,
			"country":	countryCode,
			"os":		fmt.Sprintf("%s %d", OS, ua.OS.Version.Major),
			"browser":	browser,
			"device":	device,
			"time":		time.Now().Unix(),
		})
	if err != nil {
		log.Printf("Failed to insert into database: %s", err.Error())
	}
	return nil
}

func getISOCodeFromIP(ip string) string {
	countryDB, err := geoip2.Open("geoip-city.mmdb")
	if err != nil {
		log.Printf("GeoIP2 Error: %s", err.Error())
		return ""
	}
	defer countryDB.Close()

	user, err := countryDB.Country(net.ParseIP(ip))
	if err != nil {
		log.Printf("IP Parse Error: %s", err.Error())
		return ""
	}

	return user.Country.IsoCode
}

func GetUniqueVisitors() (count int) {
	err := database.Handle.Get(&count, "SELECT COUNT(DISTINCT(uuid)) AS count FROM analytics")
	if err != nil {
		log.Printf("Could not get unique users: %s", err.Error())
		return -1
	}
	return
}

func GetCountryAnalytics() (countryAnalytics []CountryAnalytics) {
	err := database.Handle.Select(&countryAnalytics, `SELECT country,COUNT(*) AS count FROM ( SELECT uuid,country FROM analytics GROUP BY uuid ) t GROUP BY country`)
	if err != nil {
		log.Printf("Failed to get country analytics: %s", err.Error())
	}
	return
}

func GetAnalyticData() (aData []AnalyticData) {
	err := database.Handle.Select(&aData, "SELECT *,COUNT(*) AS visits FROM analytics GROUP BY uuid")
	if err != nil {
		log.Printf("Could not get analytics: %s", err.Error())
	}
	return
}

func AnalyticsByUUID(uuid string) (aData []AnalyticData) {
	err := database.Handle.Select(&aData, "SELECT * FROM analytics WHERE uuid = ?", uuid)
	if err != nil {
		log.Printf("Could not get analytics for user [%s]: %s", uuid, err.Error())
	}
	return
}

func AnalyticsByPost(title string) (aData []AnalyticData) {
	err := database.Handle.Select(&aData, "SELECT a.*,COUNT(DISTINCT a.id) AS visits FROM posts p JOIN analytics a ON REPLACE(a.page, \"/post/\", \"\") = ? GROUP BY a.uuid", title)
	if err != nil {
		log.Printf("Could not get posts: %s", err.Error())
		return nil
	}
	return
}

func PostsWithAnalytics() (posts []*AnalyticsPosts) {
	/**
	SELECT * FROM posts p JOIN analytics a
	ON a.page LIKE CONCAT('%', p.friendly_url, '%') GROUP BY p.id
	 */
	err := database.Handle.Select(&posts, "SELECT p.*,COUNT(a.uuid) AS count FROM posts p JOIN analytics a ON REPLACE(a.page, \"/post/\", \"\") = p.friendly_url GROUP BY p.id ")
	if err != nil {
		log.Printf("Could not get posts: %s", err.Error())
		return nil
	}

	formatAnalyticsPosts(posts)

	return
}

func formatAnalyticsPosts(posts []*AnalyticsPosts) {
	for _, p := range posts {
		p.Format()
	}
}