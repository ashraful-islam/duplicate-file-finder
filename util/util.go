package util

import (
	"fmt"
	"io"
	"os"
	"sort"
	"log"
	"crypto/md5"
	"github.com/ashraful-islam/duplicate-file-finder/models"
)

var (
	FILE_CHUNK_SIZE int64 = 4096
)

func ByteToStr(b []byte) string {
	return fmt.Sprintf("%x", b)
}

func GetPartialHash(filePath string, fileSize int64) (string, error) {

	var err error
	var hash string
	

	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		err = fmt.Errorf("[E_PH_01] file open failed: %v", err.Error())
		return hash, err
	}

	chunkSize := fileSize

	if fileSize > FILE_CHUNK_SIZE {
		chunkSize = FILE_CHUNK_SIZE
	}

	dataBuf := make([]byte, chunkSize)

	n, err := f.Read(dataBuf)

	if err != nil {
		err = fmt.Errorf("[E_PH_02] file read failed: %v", err.Error())
		return hash, err
	}

	if int64(n) != chunkSize {
		err = fmt.Errorf("[E_PH_03] partial read error, expected: %v but only %v was read", chunkSize, n)
		return hash, err
	}

	sum := md5.Sum(dataBuf)
	hash = ByteToStr(sum[:])
	return hash, nil
}

func GetFullHash(filePath string, fileSize int64) (string, error) {

	var err error
	var hash string
	var chunkSize int64

	// size related calculation
	numCompleteChunks := fileSize / FILE_CHUNK_SIZE
	subChunkSize := fileSize - (numCompleteChunks * FILE_CHUNK_SIZE)
	
	if numCompleteChunks > 0 {
		chunkSize = FILE_CHUNK_SIZE
	} else {
		chunkSize = subChunkSize
	}
	
	dataBuf := make([]byte, chunkSize)

	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		err = fmt.Errorf("[E_FH_01] file open failed: %v", err.Error())
		return hash, err
	}

	hasher := md5.New()

	for {

		// all done
		if numCompleteChunks == 0 && subChunkSize == 0 {
			break
		}

		// last read left or tiny file
		if numCompleteChunks == 0 && subChunkSize > 0 {
			chunkSize = subChunkSize
			subChunkSize = 0 // reset
		}

		n, err := f.Read(dataBuf)

		if err != nil {

			if err != io.EOF {
				err = fmt.Errorf("E_FH_02] file read failed: %v", err.Error())
			} else {
				// reset to nil as EOF is not proper error
				err = nil
			}

			break
		}

		// handle partial reads
		if int64(n) != chunkSize {
			err = fmt.Errorf("[E_PH_03] partial read error, expected: %v but only %v was read", chunkSize, n)
			break
		}

		// reduce counters
		if numCompleteChunks != 0 {
			numCompleteChunks--
		}

		// re calculate hash on each iteration
		hasher.Write(dataBuf)
	}
	
	// some error occured continue
	if err != nil {
		return hash, err
	}

	sum := hasher.Sum(nil)
	hash = ByteToStr(sum)
	return hash, nil
}

func SortBucketBySize(files []models.File) {
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Size < files[j].Size
	})
}

func SortBucketByFullHash(files []models.File) {
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].FullHash < files[j].FullHash
	})
}

func RemoveUniques(files []models.File) []models.File {

	count := 0
	max := len(files)
	duplicateFiles := make([]models.File, max)
	seenFiles := make(map[string]struct{}, max)

	// indices

	for i := 0; i < max - 1; i++ {

		for j := i + 1; j < max; j++ {

			if files[i].Size != files[j].Size {
				break
			}

			// no complete hash available yet
			if !files[i].HasHashes() {
				fullHash, err := GetFullHash(files[i].Path, files[i].Size)
				if err != nil {
					fmt.Printf("[E_RU01] %v", err.Error())
					break
				}
				files[i].FullHash = fullHash
			}

			if !files[j].HasHashes() {
				fullHash, err := GetFullHash(files[j].Path, files[j].Size)
				if err != nil {
					fmt.Printf("[E_RU02] %v", err.Error())
					break
				}
				files[j].FullHash = fullHash
			}

			if files[i].IsEql(files[j]) {

				// add each only if not already added
				if _, seen := seenFiles[files[i].Path]; !seen {
					seenFiles[files[i].Path] = struct{}{}
					duplicateFiles[count] = files[i]
					count++
				}

				if _, seen := seenFiles[files[j].Path]; !seen {
					duplicateFiles[count] = files[j]
					// add this file to seen files map
					seenFiles[files[j].Path] = struct{}{}
					count++
				}

			}


		}
		
	}

	// clean up
	seenFiles = nil

	// resize slice properly
	duplicateFiles = duplicateFiles[:count]
	return duplicateFiles
}

func CheckErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}