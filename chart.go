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
	Name      string
	Playcount string
	MBID      string
	URL       string
	Artist    struct {
		Name string
		MBID string
		URL  string
	}
	Image []struct {
		URL  string `json:"#text"`
		Size string
	}
	Attr struct {
		Rank string
	} `json:"@attr"`
}

type userTopAlbums struct {
	Topalbums struct {
		Album []album
		Attr  struct {
			User       string
			Page       string
			PerPage    string
			TotalPages string
			Total      string
		} `json:"@attr"`
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
	albums, loadError := load(user, period, size)
	if loadError != nil {
		printError(loadError)
		return
	}
	covers, cacheErr := cache(albums)
	if cacheErr != nil {
		printError(cacheErr)
	}
	ctx, drawErr := draw(covers, size)
	if drawErr != nil {
		printError(drawErr)
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

func cache(albums []album) ([]albumCover, error) {
	var covers []albumCover
	for i := 0; i < len(albums); i++ {
		url := albums[i].Image[3].URL
		client := &http.Client{Timeout: 10 * time.Second}
		res, err := client.Get(url)
		if err != nil {
			return covers, err
		}
		defer res.Body.Close()
		image, _, err := image.Decode(res.Body)
		if err != nil {
			return nil, err
		}
		file, err := ioutil.TempFile("/tmp", "")
		if err != nil {
			return nil, err
		}
		defer file.Close()
		path := fmt.Sprintf("%s.png", file.Name())
		if err := os.Rename(file.Name(), path); err != nil {
			return nil, err
		}
		if _, err := io.Copy(file, res.Body); err != nil {
			return nil, err
		}
		var c albumCover
		c.artist = albums[i].Artist.Name
		c.image = image
		c.path = path
		c.title = albums[i].Name
		c.width = image.Bounds().Dx()
		covers = append(covers, c)
	}
	return covers, nil
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
			ctx.DrawString(covers[i].artist, float64(x), float64(y)+h)
			ctx.DrawString(covers[i].title, float64(x), float64(y)+(h*2)+4)
			ctx.SetHexColor("ffffff")
			ctx.DrawString(covers[i].artist, float64(x)+1, float64(y)+h-1)
			ctx.DrawString(covers[i].title, float64(x)+1, float64(y)+(h*2)+3)
			i++
		}
	}
	return ctx, nil
}

func cleanup(covers []albumCover) error {
	for i := 0; i < len(covers); i++ {
		err := os.Remove(covers[i].path)
		if err != nil {
			return err
		}
	}
	return nil
}
