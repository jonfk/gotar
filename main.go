package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	flag.Parse() // get the arguments from command line

	destinationfile := flag.Arg(0)

	if destinationfile == "" {
		fmt.Println("Usage : gotar destinationfile.tar.gz source")
		os.Exit(1)
	}

	sourcedir := flag.Arg(1)

	if sourcedir == "" {
		fmt.Println("Usage : gotar destinationfile.tar.gz source-directory")
		os.Exit(1)
	}

	err := compressDirectory(sourcedir, destinationfile)
	if err != nil {
		log.Fatal(err)
	}

}

func compressDirectory(sourceDir, destinationFile string) error {

	dir, err := os.Open(sourceDir)
	if err != nil {
		return err
	}

	defer dir.Close()

	tarfile, err := os.Create(destinationFile)
	if err != nil {
		return err
	}

	defer tarfile.Close()
	var fileWriter io.WriteCloser = tarfile

	if strings.HasSuffix(destinationFile, ".gz") {
		fileWriter = gzip.NewWriter(tarfile) // add a gzip filter
		defer fileWriter.Close()             // if user add .gz in the destination filename
	}

	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return addFile(path, info, tarfileWriter)
	})
	return nil
}

func addFile(path string, fileInfo os.FileInfo, tarFileWriter *tar.Writer) error {
	if fileInfo.IsDir() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	header := new(tar.Header)
	header.Name = file.Name()
	header.Size = fileInfo.Size()
	header.Mode = int64(fileInfo.Mode())
	header.ModTime = fileInfo.ModTime()

	err = tarFileWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tarFileWriter, file)
	if err != nil {
		return err
	}

	return nil
}
