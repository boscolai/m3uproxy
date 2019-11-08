package filter

import (
	"fmt"
	"testing"

	"os"

	"github.com/stretchr/testify/require"
)

const testInputURI = "./testdata/sample.m3u"

func TestGetGroups(t *testing.T) {
	groups, err := GetGroups(testInputURI)
	require.Nil(t, err, "GetGroups returned error")
	require.NotEmpty(t, groups, "GetGroups returned empty groups")
	require.Greaterf(t, len(groups), 0, "groups length")
	for k, v := range groups {
		fmt.Printf("%s: %d\n", k, v)
	}
}

func TestFilterByGroupNames(t *testing.T) {
	groups := []string{
		"VIP USA ENTERTAINMENT",
		"USA NEWS NETWORKS",
	}

	tracks, err := FilterByGroupNames(testInputURI, groups, false)
	require.Nil(t, err, "FilterByGroupNames returned error")
	require.NotEmpty(t, tracks, "FilterByGroupNames returned empty tracks")
	require.EqualValuesf(t, 134+36, len(tracks), "number of returned tracks")

	WriteTracks(os.Stdout, tracks)
}

func TestFilterByGroupExclusion(t *testing.T) {
	groups := []string{
		"VIP USA ENTERTAINMENT",
		"USA NEWS NETWORKS",
	}

	tracks, err := FilterByGroupNames(testInputURI, groups, true)
	require.Nil(t, err, "FilterByGroupNames returned error")
	require.NotEmpty(t, tracks, "FilterByGroupNames returned empty tracks")
	require.Greaterf(t, len(tracks), 134+36, "too few tracks")

	WriteTracks(os.Stdout, tracks)
}
