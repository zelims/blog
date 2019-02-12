package models

import (
	"fmt"
	"github.com/avct/uasurfer"
	"github.com/oschwald/geoip2-golang"
	"github.com/revel/revel"
	"github.com/satori/go.uuid"
	"github.com/zelims/blog/app"
	"log"
	"net"
	"strings"
	"time"
)

type AnalyticData struct {
	CountryCode		string	`db:"country"`
	UserCount		int		`db:"count"`
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
	ipAddress := strings.Split(request.RemoteAddr, ":")[0]
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

	_, err := app.DB.NamedExec(`INSERT INTO analytics(ID,uuid,page,ip_address,country,os,browser,device,time)` +
		` VALUES(:id,:uuid,:page,:ip,:country,:os,:browser,:device,:time)`,
		map[string]interface{}{
			"id":		nil,
			"uuid":		trackID,
			"page":		request.URL.String(),
			"ip":		request.RemoteAddr,
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
	err := app.DB.Get(&count, "SELECT COUNT(DISTINCT(uuid)) AS count FROM analytics")
	if err != nil {
		log.Printf("Could not get unique users: %s", err.Error())
		return -1
	}
	return
}

func GetAnalyticData() (analyticData AnalyticData) {
	err := app.DB.Select(&analyticData, `SELECT country,COUNT(*) AS count FROM ( SELECT uuid,country FROM analytics GROUP BY uuid ) t GROUP BY country`)
	if err != nil {
		log.Printf("Failed to get analytics: %s", err.Error())
	}
	return
}