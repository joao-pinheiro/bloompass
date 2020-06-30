package bloompass

import (
	"bufio"
	"os"
	"path/filepath"
)

type addBloom func(s string) // Add to bloom filter
type console func(m string)  // log info messages to output

func ParseFiles(dirname string, fn addBloom, out console) error {

	f, err := os.Open(dirname)
	if err == nil {
		defer f.Close()
		finfo, err := f.Readdir(-1)
		if err == nil {
			for _, file := range finfo {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".txt" {
					txtfile := filepath.Join(dirname, file.Name())
					out(txtfile)
					if err = loadFile(txtfile, fn); err != nil {
						return err
					}

				}
			}
		}
	}
	return err
}

func loadFile(fname string, fn addBloom) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fn(scanner.Text())
	}

	return scanner.Err()
}
