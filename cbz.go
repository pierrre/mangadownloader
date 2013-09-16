package mangadownloader

/*
import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

func main() {
	zipFile, err := os.Create("test.zip")
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	file, err := os.Open("README.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	zipFileWriter, err := zipWriter.Create("test.md")
	if err != nil {
		panic(err)
	}

	written, err := io.Copy(zipFileWriter, file)
	if err != nil {
		panic(err)
	}
	fmt.Println(written)
}
*/
