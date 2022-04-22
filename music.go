package goacnh

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

// Song represents a K.K.Slider song as represented via the API
type Song struct {
	ID       int               `json:"id"`
	FileName string            `json:"file-name"`
	Name     map[string]string `json:"name"`
}

const (
	songNameLanguageCode string = "EUen"
	songFileExtension    string = ".mp3"
)

// SongList returns all the songs that the API provides. An error is returned if
// the request failed or a non 200 error code was returned.
func (c *Client) SongList() ([]*Song, error) {
	var songMap map[string]*Song
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetPathParam("apiVersion", strconv.Itoa(1)).
		SetResult(&songMap).
		Get("/v{apiVersion}/songs")
	if err != nil {
		return nil, fmt.Errorf("failed to request song list: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	songList := make([]*Song, 0)
	for _, value := range songMap {
		songList = append(songList, value)
	}
	return songList, nil
}

// SongByID gets a single song based on the ID provided. An error is returned if
// the request failed or a non 200 error code was returned.
func (c *Client) SongByID(id int) (*Song, error) {
	var song *Song
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetPathParam("apiVersion", strconv.Itoa(1)).
		SetPathParam("songID", strconv.Itoa(id)).
		SetResult(&song).
		Get("/v{apiVersion}/songs/{songID}")
	if err != nil {
		return nil, fmt.Errorf("failed to request song: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	return song, nil
}

// SongByName get a song based on its name. It is important to note that
// language of the name is set to EUen. An error is returned if the request
// failed or a non 200 error code was returned or no match was found.
func (c *Client) SongByName(name string) (*Song, error) {
	songList, err := c.SongList()
	if err != nil {
		return nil, err
	}
	name = strings.ToLower(name)
	for _, song := range songList {
		if strings.ToLower(song.Name[fmt.Sprintf("name-%s", songNameLanguageCode)]) == name {
			return song, nil
		}
	}
	return nil, fmt.Errorf("failed to find a match")
}

// SongDownload downloads the given track as an MP3 file to a given directory.
// The file name of the download is that specified as the file name by the API.
// The given download dir must exist before calling this. Returned is the file
// path of the download song, provided there was no error.
func (c *Client) SongDownload(song *Song, downloadDirectory string) (string, error) {
	if !dirExists(downloadDirectory) {
		return "", fmt.Errorf("destination download directory does not exist")
	}
	outputFilePath := path.Join(downloadDirectory, song.FileName) + songFileExtension
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetPathParam("apiVersion", strconv.Itoa(1)).
		SetPathParam("songID", strconv.Itoa(song.ID)).
		SetOutput(outputFilePath).
		Get("/v{apiVersion}/music/{songID}")
	if err != nil {
		return "", fmt.Errorf("failed to download background music track: %w", err)
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	return outputFilePath, nil
}

// SongDownload downloads the given track as an MP3 file to a temp directory. Th
// file name of the download is that specified as the file name by the API.
// Returned is the file path of the download song, provided there was no error.
func (c *Client) SongDownloadTemp(song *Song) (string, error) {
	return c.SongDownload(song, os.TempDir())
}
