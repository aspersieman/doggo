package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/yaml.v3"
)

type Scene struct {
	Path string `yaml:"path"`
}
type Video struct {
	OutputFile string  `yaml:"file"`
	Scenes     []Scene `yaml:"scenes"`
}
type Editor struct {
	Video Video `yml:"video"`
}

func main() {
	editor := Editor{}

	ymlFile, err := os.ReadFile("sample/sample.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal([]byte(ymlFile), &editor)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("Editor config: %+v\n", &editor)

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
	scene1 := ffmpeg.Input(editor.Video.Scenes[0].Path, ffmpeg.KwArgs{"filter_complex": "xfade=transition=slideleft:duration=3:offset=3;acrossfade=duration=2"})
	scene2 := ffmpeg.Input(editor.Video.Scenes[1].Path)
	// scene2 := ffmpeg.Input(editor.Video.Scenes[1].Path, ffmpeg.KwArgs{"filter_complex": "xfade=transition=fade:offset=1:duration=3"})
	err = ffmpeg.Concat([]*ffmpeg.Stream{
		scene1,
		scene2,
	}).Output(editor.Video.OutputFile).Run()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
