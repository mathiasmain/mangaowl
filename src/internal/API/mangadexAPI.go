package API

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Error struct {
	Id      string `json: "id"`
	Status  string `json: "status"`
	Title   string `json: "title"`
	Detail  string `json: "detail"`
	Context string `json: "context"`
}

type Description struct {
	En string `json: "en"`
}

type Title struct {
	En string `json:"en"`
}

type Name struct {
	En string `json:"en"`
}

type Attributes2 struct {
	Name Name `json:"name"`
}

type Tags struct {
	Attributes Attributes2 `json:"attributes"`
}

type Attributes struct {
	Title Title `json: "title"`
	//Description              Description `json: "description"`
	OriginalLanguage string `json: "originalLanguage"`
	Status           string `json: "status"`
	Year             uint16 `json: "year"`
	Tags             []Tags `json: "tags"`
	//LastestUploadedChapterID string      `json: "lastestUploadedChapter"`
}

type Obra struct {
	Id         string     `json: "id"`
	Type       string     `json: "type"`
	Attributes Attributes `json:"attributes"`
}

type ResponseTitlesMangaDexAPI struct {
	Result string `json: "result"`
	Data   []Obra `json: "data"`
	Limit  int    `json: "limit"`
	Offset int    `json: "offset"`
	Total  int    `json: "total"`
	Errr   Error  `json: "errors"`
}

type Attributes3 struct {
	Chapter            string `json: "chapter"`
	TranslatedLanguage string `json: "translatedLanguage"`
}

type Chapter struct {
	Id         string      `json: "id"`
	Attributes Attributes3 `json: "attributes"`
}

type ResponseChaptersMangaDexAPI struct {
	Result string    `json: "result"`
	Data   []Chapter `json: "data"`
	Limit  int       `json: "limit"`
	Offset int       `json: "offset"`
	Total  int       `json: "total"`
	Errr   Error     `json: "errors"`
}

//--------------------------------------- Implementação -----------------------------------------------
// Eu preciso definir exatamente que tipo de obra eu quero.
// Eu preciso retornar as respostas como "[]Title"
// urls:
// https://api.mangadex.org/statistics/manga/8b34f37a-0181-4f0b-8ce3-01217e9a602c
// https://mangadex.org/title/ade0306c-f4b6-4890-9edb-1ddf04df2039/solo-leveling-ragnarok

const URL_MangadexAPI = "https://api.mangadex.org/"
const URL_Mangadex = "https://mangadex.org/"

func SearchMangDexTitles(title string) (*ResponseTitlesMangaDexAPI, error) {

	var responsible ResponseTitlesMangaDexAPI
	err := requestAndJsonfyResponse(URL_MangadexAPI+"manga?title="+strings.ReplaceAll(strings.ToLower(title), " ", "-"), &responsible)
	if err != nil {
		return nil, err
	}
	return &responsible, nil
}

func GetChapterResponseByID(id string, offset int) (*ResponseChaptersMangaDexAPI, error) {
	var res ResponseChaptersMangaDexAPI // if res.Total > res.Limit {}
	search := URL_MangadexAPI + "manga/" + strings.ToLower(id) + "/feed?translatedLanguage[]=en&offset=" + strconv.Itoa(offset)
	fmt.Println("Pesquisando: ", search)
	err := requestAndJsonfyResponse(search, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func showMangasByTitle(title string) {
	res, err := SearchMangDexTitles(title)
	if err != nil {
		log.Fatalln(err)
	}

	if res.Result != "ok" {
		//if res.Errr.Status == "404" {
		//	return nil, errors.New("not found: This title does not exist.")
		//} else {
		log.Fatalln("Error:\nID: ", res.Errr.Id,
			"\nStatus: ", res.Errr.Status,
			"\nTitle: ", res.Errr.Title,
			"\nDetail: ", res.Errr.Detail)
		//}

	}

	for i, title := range res.Data {
		fmt.Println(i, title.Attributes.Title.En, title.Attributes.Tags)
	}
}

type maxNumber struct {
	numberFloat float64
	numberInt   int
	id          string
	//Bigindex    int
	//Littleindex int
}

func GetLatestChapterByID(id string) (*[2]string, error) {

	offset := 0
	ResponseArray := make([]ResponseChaptersMangaDexAPI, 0)

	var max maxNumber
	max.numberFloat = -1
	max.id = ""
	//max.Bigindex = -1
	//max.Littleindex = -1

	for {

		resultado, err := GetChapterResponseByID(id, offset)
		if err != nil {
			log.Fatalln("Olá hehe", err)
			break
		}

		if resultado.Result != "ok" {

			if resultado.Errr.Status == "404" {
				return nil, errors.New("not found: This title does not exist")
			} else {
				log.Fatalln("Error:\nID: ", resultado.Errr.Id,
					"\nStatus: ", resultado.Errr.Status,
					"\nTitle: ", resultado.Errr.Title,
					"\nDetail: ", resultado.Errr.Detail)
				break
			}

		}

		if resultado.Data[0].Id == "" {
			fmt.Println(resultado.Data[0].Id)
			log.Fatalln("I tried.Nothing returned.")
			break
		}

		ResponseArray = append(ResponseArray, *resultado)

		if (offset + resultado.Limit) > resultado.Total {
			break
		}
		offset += resultado.Limit
		time.Sleep(2 * time.Second)
	}

	// Testar com um que precise de mais de uma request para obter o objeto inteiro.
	var ChapterN float64
	ChapterN = float64(-1)
	var err error
	for _, response := range ResponseArray {
		for _, chapter := range response.Data {

			ChapterN, err = strconv.ParseFloat(chapter.Attributes.Chapter, 64)
			if err != nil {
				fmt.Println("Unsucessful convertion bettwen this string "+chapter.Attributes.Chapter+" to a int or float. Error:", err)
				return nil, errors.New("Something went wrong while trying to check the last chapter's number. Try again, please.")
			}

			if ChapterN > max.numberFloat {
				max.numberFloat = ChapterN
				max.id = chapter.Id
				//max.Bigindex = i
				//max.Littleindex = j
			}

		}
	}

	return &[2]string{strconv.FormatFloat(max.numberFloat, 'f', -1, 64), max.id}, nil
}

func Redirect(id string) {
	fmt.Println(URL_Mangadex + "title/" + id)
}

/*
[v] Comick
[v] MangaDex
[v] Bato.to
[x] MangaFire

What something.com should provide?
- Search titles.
- Add, delete, update (last readed chapter) and readAll chapters.

Which features the API should :
- Search (Comick, MangaDex, Bato)
- Bato têm menos tags, mesmo que estejam misturadas com outras coisas
- Comick tem milhões de tags, mas têm uma página para cada título.
*/
