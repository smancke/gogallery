package imglib

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

// The Sqlite filename within the library directory
var DBFilename = ".library.db"
var WriteTestFilename = ".galleryWriteTest"

type User struct {
	ID        uint       `gorm:"primary_key"`
	UserName  string     `sql:"type:varchar(50);unique_index"json:"userName"`
	NickName  string     `sql:"type:varchar(50)"json:"nickName"`
	Link      string     `sql:"type:varchar(500)"json:"link"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type Image struct {
	ID            uint       `gorm:"primary_key"`
	UserID        uint       `sql:"index"`
	User          User       `json:"user"`
	Tags          string     `json:"tags"`
	LargeFilename string     `json:"src"`
	LargeW        int        `json:"w"`
	LargeH        int        `json:"h"`
	ThumbFilename string     `json:"msrc"`
	ThumbW        int        `json:"mw"`
	ThumbH        int        `json:"mh"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     time.Time  `json:"-"`
	DeletedAt     *time.Time `json:"-"`
}

type ImageLibrary struct {
	db  *gorm.DB
	dir string
}

// Opens the library denoted by the given directory.
// If the directory does not exist, it will be created.
func (lib *ImageLibrary) Open(directoryPath string) error {
	if err := ensureWriteableDirectory(directoryPath); err != nil {
		return err
	}
	dbFilename := path.Join(directoryPath, DBFilename)
	if err := lib.openDB(dbFilename); err != nil {
		return err
	}
	lib.dir = directoryPath
	return nil
}

func (lib *ImageLibrary) openDB(filename string) error {
	log.Printf("opening sqlite3 db: %v", filename)
	gormdb, err := gorm.Open("sqlite3", filename)
	if err == nil {
		if err := gormdb.DB().Ping(); err != nil {
			log.Println("error pinging database: %v", err)
		} else {
			log.Println("can ping database")
		}

		//gormdb.LogMode(true)
		gormdb.DB().SetMaxIdleConns(2)
		gormdb.DB().SetMaxOpenConns(5)
		gormdb.SingularTable(true)

		if err := gormdb.AutoMigrate(&User{}, &Image{}).Error; err != nil {
			log.Printf("error in schema migration: %v", err)
			return err
		} else {
			log.Println("ensured db schema")
		}
	} else {
		log.Println("error opening sqlite3 db %v: %v", filename, err)
	}
	lib.db = gormdb
	return err
}

func (lib *ImageLibrary) CreateUser(user *User) error {
	return lib.db.Create(user).Error
}

func (lib *ImageLibrary) SaveUser(user *User) error {
	return lib.db.Save(user).Error
}

func (lib *ImageLibrary) UserByUsername(username string) (*User, error) {
	user := &User{}
	if err := lib.db.
		Where("user_name = ?", username).
		First(user).Error; err != nil {

		return nil, err
	}
	return user, nil
}

func (lib *ImageLibrary) CreateImage(user User, imageStream io.Reader) (*Image, error) {
	if lib == nil {
		panic("calling CreateImage on a null object")
	}

	image := &Image{}
	image.User = user
	image.UserID = user.ID

	// create image in db first to ensure the id
	if err := lib.db.Create(&image).Error; err != nil {
		return nil, err
	}

	if err := image.SaveToDirectory(lib.dir, imageStream, DefaultImageConfiguration); err != nil {
		lib.db.Delete(&image)
		return nil, err
	}

	if err := lib.db.Save(&image).Error; err != nil {
		image.DeleteFromDirectory(lib.dir)
		lib.db.Delete(&image)
		return nil, err
	}

	return image, nil
}

func (lib *ImageLibrary) GetImages() ([]*Image, error) {
	images := make([]*Image, 0, 50)
	q := lib.db.
		Model(&images).
		Preload("User").
		Order("image.created_at desc")
	if err := q.Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (lib *ImageLibrary) GetImagesByUsername(username string) ([]*Image, error) {
	user, err := lib.UserByUsername(username)
	if err != nil {
		return nil, err
	}
	images := make([]*Image, 0, 20)
	q := lib.db.
		Model(&images).
		Preload("User").
		Where("user_id = ?", user.ID).
		Order("image.created_at desc")
	if err := q.Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (lib *ImageLibrary) DeleteImage(userid, imageid uint) error {
	image := Image{}
	if err := lib.db.
		Where("user_id = ? AND id = ?", userid, imageid).
		First(&image).Error; err != nil {
		return err
	}
	errLarge, errThumb := image.DeleteFromDirectory(lib.dir)
	if errLarge != nil || errThumb != nil {
		log.Printf("error deleting image files %v, %v, %v", image, errLarge, errThumb)
	}
	return lib.db.Delete(&image).Error
}

func (lib *ImageLibrary) Close() (err error) {
	log.Printf("closing sqlite3 db")
	return lib.db.Close()
}

func ensureWriteableDirectory(dir string) error {
	dirInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		dirInfo, err = os.Stat(dir)
	}

	if err != nil || !dirInfo.IsDir() {
		return fmt.Errorf("not a directory %v", dir)
	}

	writeTest := path.Join(dir, WriteTestFilename)
	if err := ioutil.WriteFile(writeTest, []byte("writeTest"), 0644); err != nil {
		return err
	}
	if err := os.Remove(writeTest); err != nil {
		return err
	}
	return nil
}
