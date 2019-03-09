package main

import (
  "fmt"
  "flag"
  "os"
  "path/filepath"
  "github.com/ashraful-islam/duplicate-file-finder/util"
  "github.com/ashraful-islam/duplicate-file-finder/models"
)

func main() {
  var rootPath string

  // flags
  flag.StringVar(&rootPath, "path", "", "define a path to search from")
  flag.StringVar(&rootPath, "p", "", "define a path to search from (shorthand)")

  flag.Parse()

  if rootPath == "" {
    fmt.Printf("usage: dupfinder -path /home/myusername")
    os.Exit(0)
  }

  filesBucket := make(map[string][]models.File)
  var fileCount int64
  var totalSize int64

  // log start
  fmt.Println("[info] Searching for duplicates in: ", rootPath)
  fmt.Println("[status] Processing (this may take some time)")

  err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {

    if err != nil {
      fmt.Printf("Error: %v \nPath: %v\n", err, path)
      return err
    }

    // ignore irrelevant files
    if info.Name() == ".DS_Store" {
      return nil
    }

    if info.IsDir() {
      // skip git directory and contents
      if info.Name() == ".git" {
        return filepath.SkipDir
      }
      // skip directory
      return nil
    }

    // generate partial hashes  
    hash, err := util.GetPartialHash(path, info.Size())
    util.CheckErr(err)

    // if a bucket for given has does not exist, generate one
    if _, exists := filesBucket[hash]; exists != true {
      filesBucket[hash] = make([]models.File, 0)
    }

    fileEntry := models.File{ 
      Name: info.Name(),
      Size: info.Size(),
      Path: path,
      PartHash: hash,
      FullHash: "",
    }
    // add to bucket
    filesBucket[hash] = append(filesBucket[hash], fileEntry)
    fileCount++
    totalSize = totalSize + info.Size()
    
    return nil
  })
  util.CheckErr(err)

  // log progress
  fmt.Println("[status] Search Statistics: Files Scanned",fileCount," Disk Space:",totalSize)
  
  fmt.Println("[status] Checking for duplicates (this will take some time)")


  for hash, files := range filesBucket {
    
    count := len(files)
    
    if count > 1 {
      // sort by size
      util.SortBucketBySize(filesBucket[hash])
      // this will return a filtered slice
      filesBucket[hash] = util.RemoveUniques(filesBucket[hash])
      // sort by full hash so similar files are nearby to each other
      util.SortBucketByFullHash(filesBucket[hash])
    } else {
      filesBucket[hash] = nil
    }

  }

  // all done, now show results
  for _, files := range filesBucket {
    count := len(files)

    if count > 1 {
      fmt.Println("----")
      fmt.Println("Duplicate Group: ", count)
      for _, file := range files {
        fmt.Println("Name: ", file.Name)
        fmt.Println("Size: ", file.Size)
        fmt.Println("Path: ", file.Path)
        fmt.Println("Hash: ", file.FullHash)
        fmt.Println("")
      }
    }
  }

  fmt.Println("[status] Done!")
}
