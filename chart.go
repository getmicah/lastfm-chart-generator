package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fogleman/gg"
)

type album struct {
	Name   string
	Artist struct {
		Name string
	}
	Image []struct {
		URL string `json:"#text"`
	}
}

type userTopAlbums struct {
	Topalbums struct {
		Album []album
	}
}

type albumCover struct {
	artist string
	image  image.Image
	path   string
	title  string
	width  int
}

func chart(user string, period string, size int) {
	fmt.Println("Loading user...")
	albums, loadError := load(user, period, size)
	if loadError != nil {
		printError(loadError)
		return
	}
	fmt.Println("Fetching covers...")
	covers, saveErr := save(albums)
	if saveErr != nil {
		printError(saveErr)
		return
	}
	fmt.Println("Drawing chart...")
	ctx, drawErr := draw(covers, size)
	if drawErr != nil {
		printError(drawErr)
		return
	}
	filename := "collage.png"
	ctx.SavePNG(filename)
	fmt.Printf("Saved %s\n", filename)
	cleanup(covers)
}

func printError(message error) {
	fmt.Printf("Error: %v\n", message)
}

func load(user string, period string, size int) ([]album, error) {
	api := "http://ws.audioscrobbler.com/2.0/"
	key := "91289adaabc4ed1a559d5928015cd702"
	method := "user.gettopalbums"
	limit := size * size
	url := fmt.Sprintf("%s?method=%s&user=%s&period=%s&limit=%d&api_key=%s&format=json", api, method, user, period, limit, key)
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var record userTopAlbums
	if err := json.NewDecoder(res.Body).Decode(&record); err != nil {
		return nil, err
	}
	return record.Topalbums.Album, err
}

func save(albums []album) ([]albumCover, error) {
	var covers []albumCover
	for i := 0; i < len(albums); i++ {
		url := albums[i].Image[len(albums[i].Image)-1].URL
		var c albumCover
		if url == "" {
			c = blankCover(albums[i], 300)
		} else {
			var err error
			c, err = downloadCover(albums[i], url)
			if err != nil {
				return nil, err
			}
		}
		covers = append(covers, c)
	}
	return covers, nil
}

func blankCover(album album, width int) albumCover {
	var c albumCover
	img := image.NewGray(image.Rect(0, 0, width, width))
	c.artist = album.Artist.Name
	c.image = img
	c.path = ""
	c.title = album.Name
	c.width = width
	return c
}

func downloadCover(album album, url string) (albumCover, error) {
	var c albumCover
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return c, err
	}
	defer res.Body.Close()
	image, _, err := image.Decode(res.Body)
	if err != nil {
		return c, err
	}
	file, err := ioutil.TempFile("/tmp", "")
	if err != nil {
		return c, err
	}
	defer file.Close()
	path := fmt.Sprintf("%s.png", file.Name())
	if err := os.Rename(file.Name(), path); err != nil {
		return c, err
	}
	if _, err := io.Copy(file, res.Body); err != nil {
		return c, err
	}
	c.artist = album.Artist.Name
	c.image = image
	c.path = path
	c.title = album.Name
	c.width = image.Bounds().Dx()
	return c, nil
}

func draw(covers []albumCover, size int) (*gg.Context, error) {
	dx := covers[0].width
	w := dx * size
	ctx := gg.NewContext(w, w)
	i := 0
	for y := 0; y < ctx.Height(); y += dx {
		for x := 0; x < ctx.Width(); x += dx {
			if i == len(covers) {
				break
			}
			ctx.DrawImage(covers[i].image, x, y)
			if err := ctx.LoadFontFace("/Library/Fonts/Andale Mono.ttf", 14); err != nil {
				return nil, err
			}
			_, h := ctx.MeasureString(covers[i].artist)
			ctx.SetHexColor("000000")
			ctx.DrawString(covers[i].artist, float64(x)+1, float64(y)+h+2)
			ctx.DrawString(covers[i].title, float64(x)+1, float64(y)+(h*2)+6)
			ctx.SetHexColor("ffffff")
			ctx.DrawString(covers[i].artist, float64(x)+2, float64(y)+h+1)
			ctx.DrawString(covers[i].title, float64(x)+2, float64(y)+(h*2)+5)
			i++
		}
	}
	return ctx, nil
}

func cleanup(covers []albumCover) error {
	for i := 0; i < len(covers); i++ {
		if covers[i].path != "" {
			err := os.Remove(covers[i].path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
