package main

type File struct {
	fd int
}

func OpenFile(name string) (f *File, err error) {
	return nil, nil
}

func CloseFile(f *File) error {
	return nil
}

func ReadFile(f *File, offset int64, data []byte) int {
	return 0
}

func (f *File) Close() error {
	return nil
}

func (f *File) Read(offset int64, data []byte) int {
	return 0
}

var closeFile = (*File).Close
var readFile = (*File).Read

func main() {
	var data []byte
	f, _ := OpenFile("xxx")
	readFile(f, 0, data)
	closeFile(f)

	var Close = f.Close
	var Read = f.Read

	Read(0, data)
	Close()
}
