package main

import (
  "fmt"
  "flag"
  "os"
  "path/filepath"
  "github.com/ashraful-islam/duplicate-file-finder/util"
  "github.com/ashraful-islam/duplicate-file-finder/models"
)

type ProcessResult struct {
  NumScannedFiles int64
  NumDuplicateFiles int64
  SizeScannedFiles int64
  SizeDuplicateFiles int64
}

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

  // start processing
  result := Process(rootPath)

  // display small statistics
  fmt.Println("Result Statistics:")
  fmt.Println("-- Scanned:")
  fmt.Println("Files Found:", result.NumScannedFiles)
  fmt.Println("Size(bytes):", result.SizeScannedFiles)
  fmt.Println("-- Duplicates:")
  fmt.Println("Files Found:", result.NumDuplicateFiles)
  fmt.Println("Size(bytes):", result.SizeDuplicateFiles)
}

func Process(rootPath string) ProcessResult {

  filesBucket := make(map[string][]models.File)
  result := ProcessResult{}
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
    result.NumScannedFiles++
    result.SizeScannedFiles = result.SizeScannedFiles + info.Size()
    
    return nil
  })
  util.CheckErr(err)

  // log some info
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

      // duplicate files count
      result.NumDuplicateFiles += int64(count)

      // display stats
      fmt.Println("----")
      fmt.Println("Duplicate Group: ", count)
      for _, file := range files {
        fmt.Println("Name: ", file.Name)
        fmt.Println("Size: ", file.Size)
        fmt.Println("Path: ", file.Path)
        fmt.Println("Hash: ", file.FullHash)
        fmt.Println("")

        // cumulative size used by duplicate files
        result.SizeDuplicateFiles += file.Size
      }
    }
  }

  fmt.Println("[status] Done!")

  return result
}
