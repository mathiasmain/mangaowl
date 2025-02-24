package Store

import (
	"database/sql" // Package for SQL database interactions
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mathiasmain/weeblytitles/src/internal/API"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type Title struct {
	ID                string
	MangadexID        string
	Name              string
	Format            string
	Status            string
	Country           string
	Year              int
	LastReadedChapter string
	LastAPIChapter    string
	Tags              string
	AddedDate         string
}

// LastChaper's Owner
// ComickHID
// Comick title
// Comick Last Chapter
//

func NewTitleFromAPI(title *API.Obra) *Title {
	// Check if ID has a uuid lenght.
	if len(title.Id) != 36 {
		log.Fatalln("Manga ID lenght is odd, something is smelling fishy. TitleID:\n", title.Id)
	}

	if len(title.Attributes.Title.En) > 128 || len(title.Attributes.Title.En) < 1 {
		fmt.Println("The title: ", title.Attributes.Title.En, " lenght must be bettween 1 and 128 characters.")
	} // Check if

	if len(title.Type) > 32 || len(title.Type) < 2 {
		fmt.Println("The title: ", title.Attributes.Title.En, " has a format: ", title.Type, ". But its lenght must be bettween 2 and 32 characters.")
	}

	if len(title.Attributes.Status) > 32 || len(title.Attributes.Status) < 2 {
		fmt.Println("The title: ", title.Attributes.Title.En, " has a status: ", title.Attributes.Status, ". But its lenght must be bettween 2 and 32 characters.")
	}

	// Generally, the OriginalLanguage is like "ko", so it should be always small.
	if len(title.Attributes.OriginalLanguage) > 8 || len(title.Attributes.OriginalLanguage) < 1 {
		fmt.Println("The title: ", title.Attributes.Title.En, " has a Country: ", title.Attributes.OriginalLanguage, ". But its lenght must be bettween 2 and 8 characters.")
	}

	tags := ""
	for _, tag := range title.Attributes.Tags {
		tags += tag.Attributes.Name.En + ", "
	}

	if len(tags) > 192 {
		fmt.Println("The title: ", title.Attributes.Title.En, " has too many tags. Lenght: ", len(tags))
	}

	return &Title{ // The ID is only create when inserting the title.
		ID:                "",
		MangadexID:        title.Id,
		Name:              title.Attributes.Title.En,
		Format:            title.Type,
		Status:            title.Attributes.Status,
		Country:           title.Attributes.OriginalLanguage,
		Year:              int(title.Attributes.Year),
		LastReadedChapter: "",
		LastAPIChapter:    "",
		Tags:              tags[:len(tags)-2],
		AddedDate:         "",
	}

}

func NewTitleFromUser(title *Title) (*Title, error) {

	// ! The uuid from the title will only be created when inserting the title.
	// ! The LastAPIChapter value will be created from an API call before call InserTitle().
	// ! So, while there isn't a LastAPIChapter, its ok to the TitleID be empty.

	if len(title.MangadexID) != 36 {
		log.Fatalln("Something is smelling fishy.")
	}

	if len(title.Name) > 128 || len(title.Name) < 1 {
		fmt.Println("The title: ", title.Name, " lenght must be bettween 1 and 128 characters.")
	}

	if len(title.Format) > 32 || len(title.Format) < 2 {
		fmt.Println("The title: ", title.Name, " has a format: ", title.Format, ". But its lenght must be bettween 2 and 32 characters.")
	}

	if len(title.Status) > 32 || len(title.Status) < 2 {
		fmt.Println("The title: ", title.Name, " has a status: ", title.Status, ". But its lenght must be bettween 2 and 32 characters.")
	}

	if len(title.Country) > 8 || len(title.Country) < 1 {
		fmt.Println("The title: ", title.Name, " has a Country: ", title.Country, ". But its lenght must be bettween 1 and 8 characters.")
	}

	if len(title.Tags) > 192 || len(title.Tags) < 2 {
		fmt.Println("The title: ", title.Name, " need to have his tag's lenght bettween 2 and 192. Lenght: ", len(title.Tags))
	}

	if len(title.LastReadedChapter) > 8 || len(title.LastReadedChapter) < 1 {
		return nil, errors.New("incompatible lenght: Your last readed chapter is empty or too long. Please, try again. The lenght must bettween 1 and 8. Lenght: " + string(len(title.LastReadedChapter)))
	}

	//if len(title.LastAPIChapter) > 8 || len(title.LastAPIChapter) < 1 {
	//	fmt.Println("The title: ", title.Name, " last chapter hasn't a lenght bettween 2 and 192. Lenght: ", len(title.LastAPIChapter))
	//}

	return &Title{
		ID:                "",
		MangadexID:        title.MangadexID,
		Name:              title.Name,
		Format:            title.Format,
		Status:            title.Status,
		Country:           title.Country,
		Year:              title.Year,
		LastReadedChapter: title.LastReadedChapter,
		LastAPIChapter:    title.LastAPIChapter,
		Tags:              title.Tags,
		AddedDate:         "",
	}, nil

}

type LiteStore struct {
	DB *sql.DB
}

// newLiteStore initializes the SQLite database
func NewLiteStore() (*LiteStore, error) {

	db, err := sql.Open("sqlite3", "./app.db") // Open a connection to the SQLite database file named app.db
	if err != nil {
		return nil, err // Log an error and stop the program if the database can't be opened
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)

	return &LiteStore{
		DB: db,
	}, nil
}


func (ls *LiteStore) Init() error {
	if ls == nil || ls.DB == nil {
		return errors.New("database connection is not initialized")
	}
	statement :=
		"CREATE TABLE IF NOT EXISTS Titles (" +
			"TitleID varchar [36] NOT NULL PRIMARY KEY," +
			"MangadexID varchar [36] NOT NULL," +
			"Name varchar [128] NOT NULL," +
			"Format VARCHAR [32] NOT NULL," +
			"Status VARCHAR [16] NOT NULL," +
			"Origin VARCHAR [2]," +
			"Year INTEGER NOT NULL," +
			"LastReadedChapter VARCHAR[8] NOT NULL," +
			"LastAPIChapter VARCAHR[8] NOT NULL," +
			"Tags VARCHAR [192] NOT NULL," +
			"AddedDate TEXT NOT NULL);"

	pstmt, err := ls.DB.Prepare(statement)
	if err != nil {
		return err
	}

	defer pstmt.Close()

	if _, err = pstmt.Exec(); err != nil {
		return err
	}

	return nil
}


func (ls *LiteStore) InsertTitle(title *Title) error {
	if ls == nil || ls.DB == nil {
		fmt.Println("database connection is not initialized")
		return errors.New("DB: The attempt to add a title to your list failed")
	}
	stmt := "INSERT INTO Titles (TitleID, MangadexID, Name, Format, Status, Origin, Year, LastReadedChapter, LastAPIChapter,Tags,AddedDate) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	// Recebo dados brutos ou o objeto Title? <
	pstmt, err := ls.DB.Prepare(stmt)
	if err != nil {
		return err
	}

	defer pstmt.Close()

	result, err := pstmt.Exec(
		uuid.New().String(), title.MangadexID, title.Name,
		title.Format, title.Status, title.Country,
		title.Year, title.LastReadedChapter, title.LastAPIChapter, title.Tags, time.Now().String())

	if err != nil {
		return err
	}

	NumberOfRowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("An error ocurred while trying to insert a title, row not affected: ", err)
		return errors.New("failed insert: The attempt to add a title to your list failed")
	}
	if NumberOfRowsAffected != 1 {
		fmt.Println("An error ocurred while trying to insert a title: ")
		return errors.New("DB: The attempt to add a title to your list failed")
	}

	return nil
}


func (ls *LiteStore) GetAllTitles() (*[]Title, error) {
	if ls == nil || ls.DB == nil {
		fmt.Println("database connection is not initialized")
		return nil, errors.New("The attempt to add a title to your list failed")
	}
	titles := make([]Title, 0)

	pstmt, err := ls.DB.Prepare("SELECT * FROM Titles")
	if err != nil {
		return nil, err
	}

	defer pstmt.Close()

	rows, err := pstmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var title Title

		if err := rows.Scan(
			&title.ID, &title.MangadexID, &title.Name, &title.Format, &title.Status,
			&title.Country, &title.Year, &title.LastReadedChapter, &title.LastAPIChapter, &title.Tags, &title.AddedDate); err != nil {
			return nil, err
		}
		titles = append(titles, title)
	}
	return &titles, nil
}


func (ls *LiteStore) UpdateUserChapter(id string, lastChapter string) error {

	pstmt, err := ls.DB.Prepare("UPDATE Titles SET LastReadedChapter = ? WHERE TitleID = ?;")
	if err != nil {
		return err
	}

	defer pstmt.Close()

	if _, err = pstmt.Exec(lastChapter, id); err != nil {
		return err
	}
	return nil
}


func (ls *LiteStore) UpdateAPIChapter(id string, lastChapter string) error {

	pstmt, err := ls.DB.Prepare("UPDATE Titles SET LastAPIChapter = ? WHERE TitleID = ?;")
	if err != nil {
		return err
	}

	defer pstmt.Close()

	if _, err = pstmt.Exec(lastChapter, id); err != nil {
		return err
	}
	return nil
}

// Pronto!
func (ls *LiteStore) DeleteTitle(titleId string) error {

	pstmt, err := ls.DB.Prepare("DELETE FROM Titles WHERE TitleID = ?")
	if err != nil {
		return err
	}

	defer pstmt.Close()

	numRows, err := pstmt.Exec(titleId)
	if err != nil {
		return err
	}
	num, err := numRows.RowsAffected()
	if err != nil {
		return err
	}

	if num != 1 {
		return errors.New("NotFound: This title doesn't exist or isn't in your list")
	}

	return nil

}
