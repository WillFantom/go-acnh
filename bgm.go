package goacnh

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

// Weather is a weather condition that can be experienced in AC:NH
type Weather string

// BGMTrack represents a track that is played in the background of AC:NH under
// specified time and weather conditions.
type BGMTrack struct {
	ID       int     `json:"id"`
	FileName string  `json:"file-name"`
	Hour     int     `json:"hour"`
	Weather  Weather `json:"weather"`
}

const (
	SunnyWeather Weather = "Sunny"
	RainyWeather Weather = "Rainy"
	SnowyWeather Weather = "Snowy"
)

const (
	bgmMinHour       int    = 0
	bgmMaxHour       int    = 23
	bgmFileExtension string = ".mp3"
)

// BGMList returns all the background music tracks that the API provides. An
// error is returned if the request failed or a non 200 error code was returned.
func (c *Client) BGMList() ([]*BGMTrack, error) {
	var bgmMap map[string]*BGMTrack
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetPathParam("apiVersion", strconv.Itoa(1)).
		SetResult(&bgmMap).
		Get("/v{apiVersion}/backgroundmusic")
	if err != nil {
		return nil, fmt.Errorf("failed to request background music list: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	bgmList := make([]*BGMTrack, 0)
	for _, value := range bgmMap {
		bgmList = append(bgmList, value)
	}
	return bgmList, nil
}

// BGMTrackByID gets a single background music track based on the ID provided.
// An error is returned if the request failed or a non 200 error code was
// returned.
func (c *Client) BGMTrackByID(id int) (*BGMTrack, error) {
	var bgmTrack *BGMTrack
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetPathParam("apiVersion", strconv.Itoa(1)).
		SetPathParam("trackID", strconv.Itoa(id)).
		SetResult(&bgmTrack).
		Get("/v{apiVersion}/backgroundmusic/{trackID}")
	if err != nil {
		return nil, fmt.Errorf("failed to request background music track: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	return bgmTrack, nil
}

// BGMListByHour gets all the background music tracks that can be played in a
// given hour, regardless of the weather. An error is returned if the request
// failed or a non 200 error code was returned or no match was found.
func (c *Client) BGMListByHour(hour int) ([]*BGMTrack, error) {
	if hour > bgmMaxHour || hour < bgmMinHour {
		return nil, fmt.Errorf("hour must be between %d and %d", bgmMinHour, bgmMaxHour)
	}
	bgmList, err := c.BGMList()
	if err != nil {
		return nil, err
	}
	matchedList := make([]*BGMTrack, 0)
	for _, track := range bgmList {
		if track.Hour == hour {
			matchedList = append(matchedList, track)
		}
	}
	if len(matchedList) == 0 {
		return nil, fmt.Errorf("failed to find a match")
	}
	return matchedList, nil
}

// BGMListByWeather gets all the background music tracks that can be played in a
// given weather condition, regardless of the time. An error is returned if the
// request failed or a non 200 error code was returned or no match was found.
func (c *Client) BGMListByWeather(weather Weather) ([]*BGMTrack, error) {
	if weather != RainyWeather && weather != SunnyWeather && weather != SnowyWeather {
		return nil, fmt.Errorf("weather must be %s, %s, or %s", RainyWeather, SunnyWeather, SnowyWeather)
	}
	bgmList, err := c.BGMList()
	if err != nil {
		return nil, err
	}
	matchedList := make([]*BGMTrack, 0)
	for _, track := range bgmList {
		if track.Weather == weather {
			matchedList = append(matchedList, track)
		}
	}
	if len(matchedList) == 0 {
		return nil, fmt.Errorf("failed to find a match")
	}
	return matchedList, nil
}

// BGMTrackByQuery gets the background music track that can be played in a
// given weather condition, at a specified hour. An error is returned if the
// request failed or a non 200 error code was returned or no match was found.
func (c *Client) BGMTrackByQuery(hour int, weather Weather) (*BGMTrack, error) {
	if hour > bgmMaxHour || hour < bgmMinHour {
		return nil, fmt.Errorf("hour must be between %d and %d", bgmMinHour, bgmMaxHour)
	}
	if weather != RainyWeather && weather != SunnyWeather && weather != SnowyWeather {
		return nil, fmt.Errorf("weather must be %s, %s, or %s", RainyWeather, SunnyWeather, SnowyWeather)
	}
	bgmList, err := c.BGMList()
	if err != nil {
		return nil, err
	}
	for _, track := range bgmList {
		if track.Hour == hour && track.Weather == weather {
			return track, nil
		}
	}
	return nil, fmt.Errorf("failed to find a match")
}

// BGMDownload downloads the given track as an MP3 file to a given directory.
// The file name of the download is that specified as the file name by the API.
// The given download dir must exist before calling this. Returned is the file
// path of the download song, provided there was no error.
func (c *Client) BGMDownload(track *BGMTrack, downloadDirectory string) (string, error) {
	if !dirExists(downloadDirectory) {
		return "", fmt.Errorf("destination download directory does not exist")
	}
	outputFilePath := path.Join(downloadDirectory, track.FileName) + bgmFileExtension
	resp, err := c.restClient.R().
		SetHeader("Accept", "application/json").
		SetPathParam("apiVersion", strconv.Itoa(1)).
		SetPathParam("trackID", strconv.Itoa(track.ID)).
		SetOutput(outputFilePath).
		Get("/v{apiVersion}/hourly/{trackID}")
	if err != nil {
		return "", fmt.Errorf("failed to download background music track: %w", err)
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("received non-200 status code (%d)", resp.StatusCode())
	}
	return outputFilePath, nil
}

// BGMDownloadTemp downloads the given track as an MP3 file to a temp directory. Th
// file name of the download is that specified as the file name by the API.
// Returned is the file path of the download song, provided there was no error.
func (c *Client) BGMDownloadTemp(track *BGMTrack) (string, error) {
	return c.BGMDownload(track, os.TempDir())
}
