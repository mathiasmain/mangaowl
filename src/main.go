// Páginas:
// "/" -> "fetch"
// "/" -> "/userTitles"
// "/" -> "/addTitle"
// "/" -> "/removeTitle"
// "/" -> "/redirectToManga"

// "/" -> Se não estiver logado: "/login"
// "/login" -> "/"
// "/register" -> "/login"

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/joho/godotenv"
	"github.com/mathiasmain/weeblytitles/src/internal/API"
	"github.com/mathiasmain/weeblytitles/src/internal/Store"
)

func main() {
	LiteS, err := Store.NewLiteStore()
	if err != nil {
		log.Fatalln("Error while opening the DB.")
	}

	if err = LiteS.Init(); err != nil {
		log.Fatal("Error while initializing the DB.")
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error while opening .env:\n", err)
	}

	app := fiber.New()

	app.Use(logger.New())

	app.Use(limiter.New(limiter.Config{
		Max:               10,
		Expiration:        10 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Use(compress.New(compress.Config{Level: compress.LevelBestCompression /* 1 */}))

	// ------------------------- ENDPOINTS ---------------------
	app.Use("/", static.New("./src/client/dist"))

	app.Get("/api/search", SearchTitle) // search titles

	app.Get("/api/titles", func(c fiber.Ctx) error { return GetAll(LiteS, c) })

	app.Post("/api/titles", func(c fiber.Ctx) error { return CreateTitle(LiteS, c) }) // add title by

	app.Patch("/api/titles/readuntil", func(c fiber.Ctx) error { return UpdateUserChapterEndpoint(LiteS, c) }) // update user chapter from

	app.Patch("/api/titles", func(c fiber.Ctx) error { return UpdateAPIChapterEndpoint(LiteS, c) }) // update from API search

	app.Delete("/api/titles/:id", func(c fiber.Ctx) error { return DeleteHim(LiteS, c) }) // delete title by id

	port := os.Getenv("PORT")
	if port == "" {
		port = ":5450"
	}
	log.Fatal(app.Listen(port))

}

// http://127.0.0.1:5582/api/search/?title=Leveling
// TODO: Como combinar esquisa da DB e da API?
func SearchTitle(c fiber.Ctx) error {
	titleName := c.Query("title")
	if len(titleName) < 1 || len(titleName) > 200 {
		return c.Status(400).JSON(&fiber.Map{"error": "true", "message": "The name of the title is empty or ."})
	}

	res, err := API.SearchMangDexTitles(titleName)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{"error": "true", "message": "The name of the title is empty."})
	}

	if res.Result != "ok" {
		if res.Errr.Status == "404" {
			return c.Status(404).JSON(&fiber.Map{"error": "true", "message": "not found: This title does not exist."})
		} else {
			log.Fatalln("Error:\nID: ", res.Errr.Id,
				"\nStatus: ", res.Errr.Status,
				"\nTitle: ", res.Errr.Title,
				"\nDetail: ", res.Errr.Detail)
		}

	}
	titles := make([]Store.Title, 0)

	for i := 0; i < len(res.Data); i++ {
		titles = append(titles, *Store.NewTitleFromAPI(&res.Data[i]))
	}

	return c.Status(200).JSON(titles)

}

// http://127.0.0.1:5582/api/titles
func CreateTitle(LiteS *Store.LiteStore, c fiber.Ctx) error {

	var Faketitle Store.Title
	if err := c.Bind().Body(&Faketitle); err != nil {
		message := "invalid format: Please, try again. Error:" + err.Error()
		return c.Status(400).JSON(&fiber.Map{"error": "true", "message": message})
	}

	title, err := Store.NewTitleFromUser(&Faketitle)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(400).JSON(&fiber.Map{"error": "true", "message": "Something went wrong while checking the consistency of your input. Error: " + err.Error() + ". Try again, please. "})
	}

	strs, err := API.GetLatestChapterByID(title.MangadexID)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(400).JSON(&fiber.Map{"error": "true", "message": "Something went wrong while trying to check if new chapters are available. Try again, please"})
	}

	if len(strs[0]) < 1 || len(strs[0]) > 8 || len(strs[1]) != 36 {
		return c.Status(400).JSON(&fiber.Map{"error": "true", "message": "Something went wrong while checking the consistency of the new chapter data. Try again, please"})
	}
	title.LastAPIChapter = strs[0]

	if err = LiteS.InsertTitle(title); err != nil {
		fmt.Println(err.Error())
		return c.Status(400).JSON(&fiber.Map{"error": "true", "message": "Something went wrong while trying to save the title in your list. Try again, please"})
	}
	message := title.Name + " is now on your list."
	return c.Status(201).JSON(&fiber.Map{"error": "", "message": message})

}

func GetAll(LiteS *Store.LiteStore, c fiber.Ctx) error {

	titles, err := LiteS.GetAllTitles()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(&fiber.Map{"error": "Something went wrong. Try again please."})
	}

	return c.Status(200).JSON(titles)
}

func UpdateUserChapterEndpoint(LiteS *Store.LiteStore, c fiber.Ctx) error {
	mapa := c.Queries()

	ID := mapa["id"]
	LastUserChapter := mapa["chapter"]

	if len(ID) != 36 {
		return c.Status(400).JSON(&fiber.Map{"error": "true", "message": "Invalid title ID."})
	}

	if len(LastUserChapter) > 8 || len(LastUserChapter) < 1 {
		return c.Status(400).JSON(&fiber.Map{"error": "true", "message": "Invalid chapter length."})
	}

	err := LiteS.UpdateUserChapter(ID, LastUserChapter)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(&fiber.Map{"error": "Something went wrong. Try again please. " + err.Error()})
	}
	return c.Status(200).JSON(&fiber.Map{"error": "", "message": "The title's chapter was sucessfuly changed."})
}

func UpdateAPIChapterEndpoint(LiteS *Store.LiteStore, c fiber.Ctx) error {
	titles, err := LiteS.GetAllTitles()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(&fiber.Map{"error": "Something went wrong. Try again please."})
	}
	// Tome cuidado com o for loop.
	// Este evento deve ocorrer apenas em espaços de tempo específicos, sendo o menor deles, uma vez ao dia.
	for _, title := range *titles {
		// The title ID and API chapter number are suposed to be in the right format, so, i guess its fine. ;)
		strs, err := API.GetLatestChapterByID(title.MangadexID)
		if err != nil {
			fmt.Println(err.Error())
			return c.Status(500).JSON(&fiber.Map{"error": "true", "message": "Something went wrong while trying to check if new chapters are available. Try again, please"})
		}

		if len(strs[0]) < 1 || len(strs[0]) > 8 || len(strs[1]) != 36 {
			return c.Status(500).JSON(&fiber.Map{"error": "true", "message": "Something went wrong while checking the consistency of the new chapter data. Try again, please"})
		}

		if strs[0] == title.LastAPIChapter {
			time.Sleep(2 * time.Second)
			continue
		}

		err = LiteS.UpdateAPIChapter(title.ID, strs[0])
		if err != nil {
			fmt.Println(err)
			return c.Status(500).JSON(&fiber.Map{"error": "true", "message": "Something went wrong. Try again please. " + err.Error()})
		}

		time.Sleep(2 * time.Second)
	}

	return c.Status(200).JSON(&fiber.Map{"error": "", "message": "The title's chapter was sucessfuly changed."})
}

func DeleteHim(LiteS *Store.LiteStore, c fiber.Ctx) error {

	id := c.Params("id")
	if len(id) != 36 {
		c.Status(404).JSON(&fiber.Map{"error": "", "message": "The title does not exist."})
	}

	if err := LiteS.DeleteTitle(id); err != nil {
		c.Status(404).JSON(&fiber.Map{"error": "true", "message": err.Error()})
	}

	return c.Status(200).JSON(&fiber.Map{"error": "", "message": "The title was succesfully deleted."})
}
