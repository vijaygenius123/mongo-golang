package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"gopkg.in/mgo.v2"
	"log"
	"strings"
)

type star struct {
	Name      string  `json:"name"" bson:"name"`
	Photo     string  `json:"photo"" bson:"photo"`
	JobTitle  string  `json:"job_title"" bson:"job_title"`
	BirthDate string  `json:"birth_date"" bson:"birth_date"`
	Bio       string  `json:"bio"" bson:"bio"`
	TopMovies []movie `json:"top_movies"" bson:"top_movies"`
}

type movie struct {
	Title string `json:"title"" bson:"title"`
	Year  string `json:"year"" bson:"year"`
}

func main() {

	month := flag.Int("month", 1, "Month to fetch birthdays for")
	day := flag.Int("day", 1, "Day to fetch birthdays for")
	flag.Parse()
	crawl(*month, *day)
}

func crawl(month int, day int) {

	session := getMongoSession()

	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"))

	infoCollector := c.Clone()
	c.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		profileUrl := e.ChildAttr("div.lister-item-image > a", "href")
		profileUrl = e.Request.AbsoluteURL(profileUrl)
		infoCollector.Visit(profileUrl)
	})

	c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	infoCollector.OnHTML("#content-2-wide", func(e *colly.HTMLElement) {
		tmpProfile := star{}

		tmpProfile.Name = e.ChildText("h1.header > span.itemprop")
		tmpProfile.Photo = e.ChildAttr("#name-poster", "src")
		tmpProfile.JobTitle = e.ChildText("#name-job-categories > a  > span.itemprop")
		tmpProfile.BirthDate = e.ChildAttr("#name-born-info time", "datetime")
		tmpProfile.Bio = strings.TrimSpace(e.ChildText("#name-bio-text > div.name-trivia-bio-text"))
		fmt.Println(tmpProfile)

		e.ForEach("div.knownfor-title", func(_ int, kf *colly.HTMLElement) {
			tmpMovie := movie{}
			tmpMovie.Title = kf.ChildText("div.knownfor-title-role > a.knownfor-ellipsis")
			tmpMovie.Year = kf.ChildText("div.knownfor-year > span.knownfor-ellipsis")
			tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)

		})

		js, err := json.MarshalIndent(tmpProfile, "", "      ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(js))
		session.DB("mongo-golang").C("stars").Insert(tmpProfile)
	})

	url := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	c.Visit(url)
	fmt.Println(url)
}

func getMongoSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	log.Println("Connected To MongoDB")
	return s
}
