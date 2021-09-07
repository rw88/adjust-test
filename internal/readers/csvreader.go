package readers

import (
	"encoding/csv"
	"io"
	"log"
)

func CSVReader(reader io.Reader, ch chan<- []string)  {

	defer close(ch)

	csvReader := csv.NewReader(reader)
	i := 0
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else if i == 0 {
			i++
			continue
		}

		ch <- line

		i++
	}
}
