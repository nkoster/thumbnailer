package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Gebruik: programma <directory>")
		return
	}

	dir := os.Args[1]

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Fout bij het lezen van de directory:", err)
		return
	}

	toDo := []string{}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".png") && !strings.HasPrefix(file.Name(), "thumbnail_") {
			toDo = append(toDo, file.Name())
		}
	}

	for _, fileName := range toDo {
		originalFilePath := dir + "/" + fileName
		thumbnailFilePath := dir + "/thumbnail_" + fileName

		err := createThumbnail(originalFilePath, thumbnailFilePath)
		if err != nil {
			fmt.Println("Fout bij het maken van thumbnail voor", fileName, ":", err)
			continue
		}
	}
}

func createThumbnail(originalPath, thumbnailPath string) error {
	// Controleer of de thumbnail al bestaat.
	if _, err := os.Stat(thumbnailPath); err == nil {
		fmt.Println("Thumbnail bestaat al voor", thumbnailPath)
		return nil // Bestand bestaat al, dus geen actie nodig.
	} else if !os.IsNotExist(err) {
		// Er was een andere fout bij het controleren van het bestand.
		return err
	}

	file, err := os.Open(originalPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode de afbeelding.
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	thumbnail := resize.Resize(0, 80, img, resize.Lanczos3)

	// Maak een nieuw bestand voor de thumbnail.
	out, err := os.Create(thumbnailPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, thumbnail)
}
