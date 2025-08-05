package pkg

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

func parse(entities []Entity) map[string]interface{} {
	return map[string]interface{}{
		"lights":           parseLightsOn(entities),
		"speakers_playing": parseSpeakersPlaying(entities),
		"open_sensors":     parseOpenEntries(entities),
		"metrics":          parseThermometers(entities),
		"generated":        time.Now().Format("15h04m"),
	}
}

func parseLightsOn(entities []Entity) map[string]any {
	var on, off int

	for _, e := range entities {
		if !strings.HasPrefix(e.EntityID, "light.") {
			continue
		}
		if e.HasAttributeValue("is_hue_group", true) {
			continue
		}
		if e.State == "on" {
			on++
		} else {
			off++
		}
	}

	var percentOn float64
	if on+off > 0 {
		percentOn = float64(on) / float64(on+off) * 100.0
	}

	return map[string]any{
		"on":         on,
		"off":        off,
		"percent_on": percentOn,
	}
}

func parseSpeakersPlaying(entities []Entity) []string {
	var playing []string

	for _, e := range entities {
		if e.DeviceClass() == "speaker" && e.State == "playing" {
			playing = append(playing, e.FriendlyName())
		}
	}

	return playing
}

func parseOpenEntries(entities []Entity) []string {
	timeOpen := func(t string) string {
		parsed, err := time.Parse(time.RFC3339Nano, t)
		if err != nil {
			return "?h"
		}
		dur := time.Since(parsed)
		if dur > time.Hour*24 {
			return fmt.Sprintf("%dd+", int(dur.Hours()/24.0))
		}
		if dur > time.Hour {
			return fmt.Sprintf("%dh+", int(dur.Hours()))
		}
		return fmt.Sprintf("%dm", int(dur.Minutes()))
	}

	var open []string

	for _, e := range entities {
		if e.DeviceClass() == "door" && e.State == "on" {
			openDur := timeOpen(e.LastChanged)
			name := fmt.Sprintf("%s %s", openDur, e.FriendlyName())
			open = append(open, name)
		}
	}

	return open
}

func parseThermometers(entities []Entity) map[string]map[string]float64 {
	groups := map[string]map[string][]float64{
		"temperature": {
			"inside":  {},
			"outside": {},
			"garage":  {},
		},
		"humidity": {
			"inside":  {},
			"outside": {},
			"garage":  {},
		},
	}

	matchGroup := func(e Entity) (string, string, bool) {
		for _, l := range e.Labels {
			subGroups, groupOk := groups[l]
			if !groupOk {
				continue
			}
			for _, l2 := range e.Labels {
				_, subGroupOk := subGroups[l2]
				if subGroupOk {
					return l, l2, true
				}
			}
		}
		return "", "", false
	}

	ignoredMetricValues := []float64{140.0}
	for _, e := range entities {
		val, err := strconv.ParseFloat(e.State, 64)
		if err != nil {
			continue
		}

		group, subGroup, found := matchGroup(e)
		if !found {
			continue
		}

		if slices.Contains(ignoredMetricValues, val) {
			continue
		}

		groups[group][subGroup] = append(groups[group][subGroup], val)
	}

	// Reduce to float64 avg
	result := map[string]map[string]float64{}
	for metric, locs := range groups {
		result[metric] = map[string]float64{}
		for loc, values := range locs {
			if len(values) > 0 {
				sum := 0.0
				for _, v := range values {
					sum += v
				}
				result[metric][loc] = sum / float64(len(values))
			}
		}
	}

	return result
}
