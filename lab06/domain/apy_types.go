package domain

import (
	"hash/fnv"
	"strconv"
)

type Item struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Price    		int    `json:"price"`
	Quantity 		int    `json:"quantity"`
	Origin_country  string `json:"origin_country"`
}

func (b *Item) GetHash() int {
	data := fnv.New32()
 	data.Write([]byte(b.Name + strconv.Itoa(b.Price)))
 	return int(data.Sum32())
}
