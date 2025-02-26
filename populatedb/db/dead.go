package db

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed shows/*.yaml
var showFiles embed.FS

var db *gorm.DB

func init() {
	var err error
	fmt.Println("Populating database...")
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=shows port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect database: %v\n", err)
		panic("failed to connect database")
	}
	db.AutoMigrate(&Show{})
	db.AutoMigrate(&Set{})
	db.AutoMigrate(&SongPerformance{})

	files, err := showFiles.ReadDir("shows")
	if err != nil {
		fmt.Printf("Failed to read directory: %v\n", err)
		panic("failed to read directory")
	}

	for _, file := range files {
		populateDatabase(file)
	}
}

func PrintStatistics() {
	var count int64

	db.Model(&Show{}).Count(&count)
	slog.Info("Number of shows in the database", slog.Int64("count", count))

	db.Model(&Set{}).Count(&count)
	slog.Info("Number of sets in the database", slog.Int64("count", count))

	db.Model(&SongPerformance{}).Count(&count)
	slog.Info("Number of songs in the database", slog.Int64("count", count))

	db.Model(&SongPerformance{}).Distinct("title").Count(&count)
	slog.Info("Number of distinct song titles in the database", slog.Int64("count", count))

	db.Model(&Show{}).Distinct("venue").Count(&count)
	slog.Info("Number of distinct venues in the database", slog.Int64("count", count))
}

func populateDatabase(path fs.DirEntry) {
	fileContent, err := showFiles.ReadFile("shows/" + path.Name())
	if err != nil {
		log.Fatal(err)
	}

	shows := make(map[string]YamlShow)

	if err := yaml.Unmarshal(fileContent, &shows); err != nil {
		log.Fatal(err)
	}

	// var showList []Show
	for date, yamlShow := range shows {

		// there are a handlful of shows with dates like
		// 1967/01/13/0 and 1967/01/13/1 when there
		// are multiple shows in a day.  We only want
		// the date and this code truncates the trailing /[number].
		if count := strings.Count(date, "/"); count > 2 {
			parts := strings.SplitN(date, "/", 4)
			date = fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
		}

		parsedDate, err := time.Parse("2006/01/02", date)
		if err != nil {
			log.Fatalf("Failed to parse date: %v", err)
		}

		show := &Show{
			Date:  parsedDate,
			Venue: yamlShow.Venue,
			City:  yamlShow.City,
			State: yamlShow.State,
			Sets:  make([]Set, 0),
		}
		for setCounter, set := range yamlShow.Setlist {
			songs := getSetList(set)
			songs.SetNumber = setCounter + 1
			show.Sets = append(show.Sets, songs)
		}
		db.Create(&show)
	}
}

func getSetList(set map[string]any) Set {
	theSet := Set{}
	for songCounter, song := range set[":songs"].([]interface{}) {
		songName := song.(map[string]interface{})[":name"].(string)
		newSong := SongPerformance{Title: songName, OrderInSet: songCounter + 1}
		theSet.SongPerformances = append(theSet.SongPerformances, newSong)
	}
	return theSet
}
