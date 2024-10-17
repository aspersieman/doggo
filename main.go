package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/yaml.v3"
)

type (
	FilterType     string
	FilterCategory string
)

const (
	FILTER_TYPE_COMPLEX FilterType = "filter_complex"

	FILTER_CATEGORY_SLIDE_LEFT  FilterCategory = "slideleft"
	FILTER_CATEGORY_SLIDE_RIGHT FilterCategory = "slideright"
	FILTER_CATEGORY_SLIDE_UP    FilterCategory = "slideup"
	FILTER_CATEGORY_SLIDE_DOWN  FilterCategory = "slidedown"

	FILTER_CATEGORY_SMOOTH_LEFT  FilterCategory = "smoothleft"
	FILTER_CATEGORY_SMOOTH_RIGHT FilterCategory = "smoothright"
	FILTER_CATEGORY_SMOOTH_UP    FilterCategory = "smoothup"
	FILTER_CATEGORY_SMOOTH_DOWN  FilterCategory = "smoothdown"
)

// List of filters: https://trac.ffmpeg.org/wiki/Xfade
type Filter struct {
	Type     FilterType     `yaml:"type"`
	Category FilterCategory `yaml:"category"`
	Duration int            `yaml:"duration"`
	Offset   int            `yaml:"offset"`
}
type Scene struct {
	Path   string `yaml:"path"`
	Filter Filter `yaml:"filter"`
}
type Video struct {
	OutputFile string  `yaml:"file"`
	Scenes     []Scene `yaml:"scenes"`
}
type Editor struct {
	Video Video `yml:"video"`
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Println("error: Argument required: yaml file name")
		usage()
		log.Fatalln("exiting")
	}
	arg := args[0]
	if arg == "help" {
		usage()
		os.Exit(3)
	}
	ymlFileName := arg

	editor := Editor{}

	ymlFile, err := os.ReadFile(ymlFileName)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal([]byte(ymlFile), &editor)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Editor running...\n")
	err = os.Remove(editor.Video.OutputFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("error deleting file: %v", err)
	}
	outputPath := filepath.Dir(editor.Video.OutputFile)
	fmt.Printf("Creating storage folder: %s\n", outputPath)
	err = os.MkdirAll(outputPath, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("error creating directory: %v", err)
	}
	fmt.Printf("  Output: %s\n", editor.Video.OutputFile)
	scenes := make([]*ffmpeg.Stream, 0)
	for _, s := range editor.Video.Scenes {
		params, err := GetFilterParams(s.Filter)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		scene := ffmpeg.Input(s.Path, *params)
		scenes = append(scenes, scene)
	}
	// scene1 := ffmpeg.Input(editor.Video.Scenes[0].Path, ffmpeg.KwArgs{"filter_complex": "xfade=transition=slideleft:duration=3:offset=3;acrossfade=duration=2"})
	// scene2 := ffmpeg.Input(editor.Video.Scenes[1].Path)
	err = ffmpeg.Concat(scenes).Output(editor.Video.OutputFile).Run()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func GetFilterParams(filter Filter) (*ffmpeg.KwArgs, error) {
	if filter.Type == "filter_complex" {
		val := ""
		if filter.Category != "" {
			val = fmt.Sprintf("xfade=transition=%s:duration=%d:offset=%d;acrossfade=duration=2", filter.Category, filter.Duration, filter.Offset)
		}
		return &ffmpeg.KwArgs{"filter_complex": val}, nil
	}
	if filter.Type != "" {
		fmt.Printf("warning: filter type '%s' not recognised\n", filter.Type)
	}
	return &ffmpeg.KwArgs{}, nil
}

func usage() {
	bin := os.Args[0]
	fmt.Printf("\nUsage: %s [option]\n\nOptions:\n\tfile.yml\tCreate a video from the instructions in the yaml file\n\thelp\t\tShow this help\n\n", bin)
}
