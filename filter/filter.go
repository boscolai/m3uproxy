package filter

import (
	"strings"

	"fmt"
	"io"

	"github.com/jamesnetherton/m3u"
)

const tagGroupTitle = "group-title"

func FilterByGroupNames(uri string, groupList []string) ([]m3u.Track, error) {
	groupMap := make(map[string]bool)
	for _, g := range groupList {
		groupMap[g] = true
	}
	var selected []m3u.Track
	playlist, err := m3u.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input m3u: %s", err)
	}
	for _, track := range playlist.Tracks {
		for _, tag := range track.Tags {
			if strings.ToLower(tag.Name) == tagGroupTitle {
				if _, inSelectedGroups := groupMap[tag.Value]; inSelectedGroups {
					selected = append(selected, track)
				}
			}
		}
	}
	return selected, nil
}

func GetGroups(uri string) (map[string]int, error) {
	playlist, err := m3u.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input m3u: %s", err)
	}
	groups := map[string]int{}
	for _, track := range playlist.Tracks {
		for _, tag := range track.Tags {
			if strings.ToLower(tag.Name) == tagGroupTitle {
				count, exists := groups[tag.Value]
				if !exists {
					count = 0
				}
				groups[tag.Value] = count + 1
			}
		}
	}
	return groups, nil
}

func WriteTrack(w io.Writer, track m3u.Track) error {
	if _, err := fmt.Fprintf(w, "#EXTINF:%d", track.Length); err != nil {
		return fmt.Errorf("failed to write track data: %s", err)
	}
	for _, tag := range track.Tags {
		if _, err := fmt.Fprintf(w, ` %s="%s"`, tag.Name, tag.Value); err != nil {
			return fmt.Errorf("failed to write track data: %s", err)
		}
	}
	if _, err := fmt.Fprintf(w, ",%s\n%s\n", track.Name, track.URI); err != nil {
		return fmt.Errorf("failed to write track data: %s", err)
	}
	return nil
}

func WriteTracks(w io.Writer, tracks []m3u.Track) error {
	if _, err := fmt.Fprintln(w, "#EXTM3U"); err != nil {
		return fmt.Errorf("failed to write m3u header: %s", err)
	}
	for _, track := range tracks {
		if err := WriteTrack(w, track); err != nil {
			return fmt.Errorf("failed to write track: %s: %s", track.Name, err)
		}
	}
	return nil
}
