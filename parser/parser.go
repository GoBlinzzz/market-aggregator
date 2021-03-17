package parser

import (
	"golang.org/x/net/html"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func Search(text string, query string) []*Item { //gets items from sources and sorts them by query params
	text = url.QueryEscape(text)

	itemsW := getItems("https://wildberries.ru", "https://www.wildberries.ru/catalog/0/search.aspx?xsearch=true&search="+text, params1)
	shortenPriceWB(itemsW)
	itemsC := getItems("https://citilink.ru", "https://www.citilink.ru/search/?text="+text, params2)
	itemsE := getItems("https://www.eldorado.ru", "https://www.eldorado.ru/search/catalog.php?q="+text, params3)
	var items []*Item
	for i := 0; i < len(itemsW) || i < len(itemsC) || i < len(itemsE); i++ {
		if i < len(itemsW) {
			items = append(items, itemsW[i])
		}
		if i < len(itemsC) {
			items = append(items, itemsC[i])
		}
		if i < len(itemsE) {
			items = append(items, itemsE[i])
		}
	}
	convertPriceStringToInt(items)
	sortItems(query, items)
	return items
}

func convertPriceStringToInt(items []*Item) {
	for _, item := range items {
		if item.Price == "" {
			item.Price = "Нет в наличии"
			item.intPrice = math.MaxInt16
		} else {
			runes := []byte(item.Price)
			strInt := 0
			for _, r := range runes {
				if r >= '0' && r <= '9' {
					strInt = 10*strInt + int(r) - '0'
				}
			}
			item.intPrice = strInt
		}
	}
}

func shortenPriceWB(items []*Item) {
	for _, item := range items {
		runes := []rune(item.Price)
		for i := 0; i < len(runes); i++ {
			if runes[i] == '₽' {
				runes = runes[:i+1]
				i = len(runes)
			}
		}
		item.Price = string(runes)
	}
}

func sortItems(query string, items []*Item) {
	switch query {
	case "aprice":
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].intPrice < items[j].intPrice
		})
	case "dprice":
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].intPrice > items[j].intPrice
		})
	case "rating":
		sort.SliceStable(items, func(i, j int) bool {
			var (
				a = float32(items[i].Rating.Count)
				b = float32(items[j].Rating.Count)
			)
			if items[i].Rating.WithHalf {
				a += 0.5
			}
			if items[j].Rating.WithHalf {
				b += 0.5
			}
			if a == b {
				return items[i].ReviewCount > items[j].ReviewCount
			}
			return a > b
		})
	}
}

func getHTMLNode(link string) (*html.Node, bool) { //request to url, if it's OK returns html.body
	reqURL, _ := url.Parse(link)

	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36"},
			"Accept":     {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		},
	}

	if response, err := http.DefaultClient.Do(req); err != nil {
		return nil, false
	} else {
		defer response.Body.Close()
		if response.StatusCode == http.StatusOK {
			if node, err := html.Parse(response.Body); err != nil {
				return nil, false
			} else {
				return node, true
			}
		}
	}
	return nil, false
}

func getAttr(node *html.Node, key string) string { //returns value of key attr
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func matchNode(node *html.Node, tag, key, value string) bool { //checks if the node matches to the template
	if key == "" {
		return node != nil && node.Type == html.ElementNode && node.Data == tag
	} else {
		return node != nil && node.Type == html.ElementNode && node.Data == tag && getAttr(node, key) == value
	}
}

func searchInNode(node *html.Node, params *pattern) *html.Node { //recursively searches in node for appropriate node
	if node == nil {
		return nil
	}
	if matchNode(node, params.tag, params.key, params.value) {
		if params.num > 1 {
			params.num--
		} else {
			return node
		}
	} else {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if elem := searchInNode(child, params); elem != nil {
				return elem
			}
		}
	}
	return nil
}

func getTextFromNode(node *html.Node) (text string) {
	if node == nil {
		return
	}
	if node.Type == html.TextNode {
		text += strings.TrimSpace(node.Data)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		text += getTextFromNode(child)
	}
	return
}

func getItems(host, link string, params [9][]pattern) (items []*Item) { //parses website by searching items of products
	if node, status := getHTMLNode(link); !status {
		return nil
	} else {
		for _, paramsIter := range params[0] {
			node = searchInNode(node, &paramsIter)
		}
		if node == nil {
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if itemCandidate := searchInNode(child, &params[1][0]); itemCandidate != nil {
				items = append(items, parseItem(itemCandidate, params))
			}
		}
	}
	for _, item := range items { //adds sourceMarket and corrects links
		switch host {
		case "https://wildberries.ru":
			item.SourceMarket = "wb"
		case "https://citilink.ru":
			item.SourceMarket = "ctl"
		case "https://www.eldorado.ru":
			item.SourceMarket = "eld"
		}
		item.Link = host + item.Link
		item.ImageSrc = strings.TrimSpace(item.ImageSrc)
		if !strings.HasPrefix(item.ImageSrc, "https:") {
			item.ImageSrc = "https:" + item.ImageSrc
		}
	}
	return
}

func parseItem(node *html.Node, params [9][]pattern) *Item {
	item := &Item{}
	for _, p := range params[2] {
		node = searchInNode(node, &p)
	}
	item.Link = getLinkFromItem(node, params[3])
	item.ImageSrc = getImageFromItem(node, params[4])
	item.Name = getInfoFromItem(node, params[5])
	item.Price = getInfoFromItem(node, params[6])
	item.Rating = getRatingFromItem(node, params[7])
	if item.Rating == nil {
		return item
	}
	reviewsAmount := getInfoFromItem(node, params[8])
	for i := len(reviewsAmount) - 1; i > 0; i-- {
		if reviewsAmount[i] < '0' || reviewsAmount[i] > '9' {
			reviewsAmount = reviewsAmount[:i]
		}
	}
	item.ReviewCount, _ = strconv.Atoi(reviewsAmount)
	return item
}

func getLinkFromItem(node *html.Node, params []pattern) string {
	for _, p := range params {
		node = searchInNode(node, &p)
	}
	return getAttr(node, "href")
}

func getImageFromItem(node *html.Node, params []pattern) (attr string) {
	for _, p := range params {
		node = searchInNode(node, &p)
	}
	attr = getAttr(node, "src")
	if strings.Contains(attr, ".gif") {
		attr = "https:" + getAttr(node, "data-original")
	}
	return
}

func getInfoFromItem(node *html.Node, params []pattern) string {
	for _, p := range params {
		node = searchInNode(node, &p)
	}
	return getTextFromNode(node)
}

func getRatingFromItem(node *html.Node, params []pattern) *Stars {
	for _, p := range params {
		node = searchInNode(node, &p)
	}
	if node == nil {
		return &Stars{Count: 0, WithHalf: false}
	} else if node.FirstChild == nil {
		class := getAttr(node, "class")
		return &Stars{Count: int(class[len(class)-1]) - '0', WithHalf: false}
	} else if getAttr(node, "class") == "tevqf5-0 cbJQML" {
		count := 0
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if strings.HasSuffix(getAttr(child, "class"), "r") {
				count++
			}
		}
		return &Stars{Count: count, WithHalf: false}
	}
	floatStr := getTextFromNode(node)
	rating, _ := strconv.ParseFloat(floatStr, 10)
	if rating-float64(int(rating)) >= 0.5 {
		return &Stars{Count: int(rating), WithHalf: true}
	}
	return &Stars{Count: int(rating), WithHalf: false}
}
