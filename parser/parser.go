package parser

import (
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
)

type Item struct {
	Ref string
	ImageSrc string `json:"image"`
	Name string `json:"title"`
	Price string `json:"price"`
}

type pattern struct {
	tag, key, value string
	num int
}

type TemplateJSON struct {
	Count int `json:"count"`
	Items []*Item `json:"items"`
}

//params    0 - path til items
//			1 - cycle by items
//			2 - parsing items
//			3 - link
//			4 - img
//			5 - name
//			6 - price
var (
	params1 = [7][]pattern{
	{{"div", "id", "body-layout", 1}, {"div", "class", "left-bg", 1}, {"div", "class", "trunkOld", 1}, {"div", "id", "catalog", 1}, {"div", "id", "catalog-content", 1}, {"div", "class", "catalog_main_table j-products-container", 1}},
	{{"div", "class", "dtList i-dtList j-card-item ", 1}},
	{{"div", "class", "dtList-inner", 1}, {"span", "itemtype", "http://schema.org/Thing", 1}, {"span", "itemtype", "http://schema.org/CreativeWork", 1}, {"span", "itemtype", "http://schema.org/MediaObject", 1}, {"a", "class", "ref_goods_n_p j-open-full-product-card", 1}},
	{},
	{{"div", "class", "l_class", 1}, {"img", "class", "thumbnail", 1}},
	{{"div", "class", "dtlist-inner-brand", 1}, {"div", "class", "dtlist-inner-brand-name", 1}},
	{{"div", "class", "dtlist-inner-brand", 1}, {"div", "class", "j-cataloger-price", 1}, {"span", "class", "price", 1}}}
	params2 = [7][]pattern {
	{{"div", "class", "MainWrapper", 1}, {"div", "class", "MainLayout js--MainLayout", 1}, {"main", "class", "MainLayout__main js--MainLayout__main", 1}, {"div", "class", "SearchResults", 1}, {"div", "class", "Container Container_has-grid ProductCardCategoryList js--ProductCardCategoryList SearchResults__product-container", 1}, {"div", "class", "block_data__gtm-js block_data__pageevents-js listing_block_data__pageevents-js ProductCardCategoryList__products-container", 1}, {"div", "class", "ProductCardCategoryList__grid-container", 1}, {"div", "class", "ProductCardCategoryList__grid", 1}, {"section", "class", "GroupGrid js--GroupGrid GroupGrid_has-column-gap GroupGrid_has-row-gap", 1}},
	{{"div", "class", "product_data__gtm-js product_data__pageevents-js  ProductCardVertical js--ProductCardInListing ProductCardVertical_normal ProductCardVertical_shadow-hover ProductCardVertical_separated", 1}},
	{},
	{{"a", "class", " ProductCardVertical__link link_gtm-js  Link js--Link Link_type_default", 1}},
	{{"div", "class", "ProductCardVerticalLayout ProductCardVertical__layout", 1}, {"div", "class", "ProductCardVerticalLayout__header", 1}, {"div", "class", "ProductCardVerticalLayout__wrapper-cover-image ProductCardVertical__layout-cover-image", 1}, {"a", "class", "ProductCardVertical__image-link  Link js--Link Link_type_default", 1}, {"div", "class", "ProductCardVertical__image-wrapper", 1}, {"div", "class", "ProductCardVertical__picture-container ", 1}, {"img", "class", "ProductCardVertical__picture js--ProductCardInListing__picture", 1}},
	{{"div", "class", "ProductCardVerticalLayout ProductCardVertical__layout", 1}, {"div", "class", "ProductCardVerticalLayout__header", 1}, {"div", "class", "ProductCardVerticalLayout__wrapper-description ProductCardVertical__layout-description", 1}, {"div", "class", "ProductCardVertical__description ", 1}, {"a", "class", " ProductCardVertical__name  Link js--Link Link_type_default", 1}},
	{{"div", "class", "ProductCardVerticalLayout ProductCardVertical__layout", 1}, {"div", "class", "ProductCardVerticalLayout__footer", 1}, {"div", "class", "ProductCardVerticalLayout__wrapper-price", 1}, {"div", "class", "ProductCardVertical_mobile", 1}, {"div", "class", "ProductCardVertical__price-with-amount", 1}, {"div", "class", "ProductPrice ProductPrice_default ProductCardVerticalPrice__price-current", 1}},}
)


func Search(text string) []*Item {
	text = url.QueryEscape(text)
	return append(getItems("https://wildberries.ru", "https://www.wildberries.ru/catalog/0/search.aspx?xsearch=true&search=" + text, params1), getItems("https://citilink.ru", "https://www.citilink.ru/search/?text=" + text, params2)...)
}

func getHTMLNode (link string) (*html.Node, bool) { //request to url, if it's OK returns html.body
	reqURL, _ := url.Parse(link)

	req := &http.Request{
		Method: "GET",
		URL: reqURL,
		Header: map[string][]string {
			"User-Agent" : {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36"},
			"Accept" : {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
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

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func matchNode(node *html.Node, tag, key, value string) bool {
	if key == "" {
		return node != nil && node.Type == html.ElementNode && node.Data == tag
	} else {
		return node != nil && node.Type == html.ElementNode && node.Data == tag && getAttr(node, key) == value
	}
}

func searchInNode(node *html.Node, params *pattern) *html.Node {
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
	if node.Type == html.TextNode {
		text += strings.TrimSpace(node.Data)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		text += getTextFromNode(child)
	}
	return
}

func getItems(host, link string, params [7][]pattern) (items []*Item) { //parses website by searching items of products
	if node, status := getHTMLNode(link); !status {
		return nil
	} else {
		for _, paramsIter := range params[0] {
			node = searchInNode(node, &paramsIter)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if itemCandidate := searchInNode(child, &params[1][0]); itemCandidate != nil {
				items = append(items, parseItem(itemCandidate, params))
			}
		}
	}
	for _, item := range items {
		item.Ref = host + item.Ref
		item.ImageSrc = strings.TrimSpace(item.ImageSrc)
		if !strings.HasPrefix(item.ImageSrc, "https:") {
			item.ImageSrc = "https:" + item.ImageSrc
		}
	}
	return
}

func parseItem(node *html.Node, params [7][]pattern) *Item {
	item := &Item{}
	for _, p := range params[2] {
		node = searchInNode(node, &p)
	}
	item.Ref = getLinkFromItem(node, params[3])
	item.ImageSrc = getImageFromItem(node, params[4])
	item.Name = getInfoFromItem(node, params[5])
	item.Price = getInfoFromItem(node, params[6])
	return item
}

func getLinkFromItem(node *html.Node, params []pattern) string {
	for _, p := range params {
		node = searchInNode(node, &p)
	}
	return getAttr(node, "href")
}

func getImageFromItem(node *html.Node, params []pattern) string {
	for _, p := range params {
		node = searchInNode(node, &p)
	}
	return getAttr(node, "src")
}

func getInfoFromItem(node *html.Node, params []pattern) string {
	for _, p := range params {
		node = searchInNode(node, &p)
	}
	return getTextFromNode(node)
}
