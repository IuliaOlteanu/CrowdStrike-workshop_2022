package domain

import (
	"hash/fnv"
	"strconv"
)

type Book struct {
	Title               string   `json:"title"`
	Author              string   `json:"author"`
	Publication_year    int      `json:"publication_year"`
	Number_of_downloads int      `json:"number_of_downloads"`
	Tag_list            []string `json:"tag_list"`
}

func (b *Book) GetHash() int {
	data := fnv.New32()
	data.Write([]byte(b.Title + b.Author + strconv.Itoa(b.Publication_year)))
	return int(data.Sum32())
}
