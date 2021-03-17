package cart

import (
	"encoding/json"
	"io/ioutil"
	"market-backend/parser"
	"os"
)

type Cart struct {
	Id string `json:"id"`
	Items []*parser.Item `json:"items"`
}

func AddToCart(key string, newItem *parser.Item) {
	carts := readCartsFromStorage()
	for _, c := range carts {
		if c.Id == key {
			c.Items = append(c.Items,newItem)
			return
		}
	}
	carts = append(carts, &Cart{Id: key, Items: []*parser.Item{newItem}})
	writeCartsToStorage(carts)
}

func GetCart(key string) []byte {
	carts := readCartsFromStorage()
	for _, c := range carts {
		if c.Id == key {
			items, _ := json.Marshal(&parser.TemplateJSON{Count: len(c.Items), Items: c.Items})
			return items
		}
	}
	items, _ := json.Marshal(&parser.TemplateJSON{Count: 0, Items: nil})
	return items
}

func readCartsFromStorage() (carts []*Cart) {
	jsonStorage, _ := os.Open("cart.json")

	defer jsonStorage.Close()

	data, _ := ioutil.ReadAll(jsonStorage)

	_ = json.Unmarshal(data, &carts)
	return
}

func writeCartsToStorage(carts []*Cart) {
	data, _ := json.Marshal(carts)

	_ = ioutil.WriteFile("cart.json", data, 0644)
}