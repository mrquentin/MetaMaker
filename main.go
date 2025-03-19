package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
) //TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click

// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Error getting current directory:", err)
	}

	files, err := os.ReadDir(currentDir)
	if err != nil {
		log.Fatalln("Error reading directory:", err)
	}

	mcMetaFile, err := getMcMeta(files)
	if err != nil {
		log.Fatalln("Error when searching for McMeta:", err)
	}

	log.Default().Println("McMetaFile:", mcMetaFile)

	pngs, err := filterPNGs(files)
	if err != nil {
		log.Fatalln("could not filter PNGs:", err)
	}

	pngToIgnore := strings.TrimSuffix(mcMetaFile.Name(), ".mcmeta")
	for _, png := range pngs {
		if png.Name() != pngToIgnore {
			log.Default().Println("Processing:", png.Name())
			newFileName := png.Name() + ".mcmeta"
			err := copyFile(mcMetaFile.Name(), newFileName)
			if err != nil {
				log.Fatalln("Error copying file:", err)
			}
		}
	}

}

func filterPNGs(files []os.DirEntry) ([]os.DirEntry, error) {
	var pngs []os.DirEntry
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".png" {
			pngs = append(pngs, file)
		}
	}
	return pngs, nil
}

func getMcMeta(files []os.DirEntry) (os.DirEntry, error) {
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".mcmeta" {
			return file, nil
		}
	}
	return nil, fmt.Errorf("no .mcmeta file found")
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
