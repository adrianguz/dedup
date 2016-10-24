package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var fileList map[int64][]string
var dupeList map[string][]string
var path = flag.String("path", "", "Set path to check")

func WalkFunc(path string, info os.FileInfo, err error) error {
	if _, ok := fileList[info.Size()]; ok {
		fileList[info.Size()] = append(fileList[info.Size()], path)
	} else {
		fileList[info.Size()] = []string{path}
	}
	fmt.Println(path)
	return nil
}

func crawl(path string) {
	fmt.Printf("\nWalking %q \n", path)
	filepath.Walk(path, WalkFunc)
}

func hashFile(filename string) (string, error) {
	var checkSum []byte
	hash := md5.New()
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return string(hash.Sum(checkSum)), nil
}

func generateDup() error {
	for _, v := range fileList {
		if len(v) > 1 {
			for _, n := range v {
				hash, err := hashFile(n)
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
				dupeList[hash] = append(dupeList[hash], n)
			}
		}
	}
	for k, v := range dupeList {
		if len(v) < 2 {
			delete(dupeList, k)
		}
	}
	return nil
}

func eliminateDup() {
	for _, v := range dupeList {
		master := v[0]
		slaves := v[1:]
		for _, i := range slaves {
			fmt.Println(i)
			os.Remove(i)
			err := os.Link(master, i)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func main() {
	flag.Parse()
	if *path == "" {
		fmt.Println("--path flag must be specified.")
	}

	fileList = make(map[int64][]string)
	dupeList = make(map[string][]string)

	crawl(*path)
	generateDup()

	if dupeList != nil {
		fmt.Println("\nDuplicates found:")
		for k, v := range dupeList {
			for _, n := range v {
				fmt.Println("Duplicate", k, n)
			}
		}
		eliminateDup()
	} else {
		fmt.Println("No duplicates found")
	}
}
