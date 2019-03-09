package models

type File struct {
	Name string
	Size int64
	Path string
	PartHash string
	FullHash string
}

func (file *File) IsEql(anotherFile File) bool {
	return file.Size == anotherFile.Size && file.FullHash == anotherFile.FullHash
}

func (f *File) HasHashes() bool {
	return f.PartHash != "" && f.FullHash != ""
}