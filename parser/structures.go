package parser

type Item struct {
	Link         string `json:"link"`
	SourceMarket string `json:"sourceMarket"`
	ImageSrc     string `json:"image"`
	Name         string `json:"title"`
	Price        string `json:"price"`
	Rating       *Stars `json:"rating"`
	ReviewCount  int    `json:"reviewCount"`
	intPrice     int
}

type pattern struct {
	tag, key, value string
	num             int
}

type Stars struct {
	Count int
	WithHalf bool
}

type TemplateJSON struct {
	Count int     `json:"count"`
	Items []*Item `json:"items"`
}

//params    0 - path til items
//			1 - cycle by items
//			2 - parsing items
//			3 - link
//			4 - img
//			5 - name
//			6 - price
//			7 - rating
//			8 - reviews
var (
	params1 = [9][]pattern{ //wildberries tag structure
		{{"div", "id", "body-layout", 1}, {"div", "class", "left-bg", 1}, {"div", "class", "trunkOld", 1}, {"div", "id", "catalog", 1}, {"div", "id", "catalog-content", 1}, {"div", "class", "catalog_main_table j-products-container", 1}},
		{{"div", "class", "dtList i-dtList j-card-item ", 1}},
		{{"div", "class", "dtList-inner", 1}, {"span", "itemtype", "http://schema.org/Thing", 1}, {"span", "itemtype", "http://schema.org/CreativeWork", 1}, {"span", "itemtype", "http://schema.org/MediaObject", 1}, {"a", "class", "ref_goods_n_p j-open-full-product-card", 1}},
		{},
		{{"div", "class", "l_class", 1}, {"img", "class", "thumbnail", 1}},
		{{"div", "class", "dtlist-inner-brand", 1}, {"div", "class", "dtlist-inner-brand-name", 1}},
		{{"div", "class", "dtlist-inner-brand", 1}, {"div", "class", "j-cataloger-price", 1}, {"span", "class", "price", 1}},
		{{"span", "itemtype", "http://schema.org/Intangible", 1}, {"span", "itemtype", "http://schema.org/Rating", 1}, {"span", "itemprop", "aggregateRating", 1}},
		{{"span", "itemtype", "http://schema.org/Intangible", 1}, {"span", "itemtype", "http://schema.org/Rating", 1}, {"span", "class", "dtList-comments-count c-text-sm", 1}},
	}
	params2 = [9][]pattern{ //citilink tag structure
		{{"div", "class", "MainWrapper", 1}, {"div", "class", "MainLayout js--MainLayout", 1}, {"main", "class", "MainLayout__main js--MainLayout__main", 1}, {"div", "class", "SearchResults", 1}, {"div", "class", "Container Container_has-grid ProductCardCategoryList js--ProductCardCategoryList SearchResults__product-container", 1}, {"div", "class", "block_data__gtm-js block_data__pageevents-js listing_block_data__pageevents-js ProductCardCategoryList__products-container", 1}, {"div", "class", "ProductCardCategoryList__grid-container", 1}, {"div", "class", "ProductCardCategoryList__grid", 1}, {"section", "class", "GroupGrid js--GroupGrid GroupGrid_has-column-gap GroupGrid_has-row-gap", 1}},
		{{"div", "class", "product_data__gtm-js product_data__pageevents-js  ProductCardVertical js--ProductCardInListing ProductCardVertical_normal ProductCardVertical_shadow-hover ProductCardVertical_separated", 1}},
		{},
		{{"a", "class", " ProductCardVertical__link link_gtm-js  Link js--Link Link_type_default", 1}},
		{{"div", "class", "ProductCardVerticalLayout ProductCardVertical__layout", 1}, {"div", "class", "ProductCardVerticalLayout__header", 1}, {"div", "class", "ProductCardVerticalLayout__wrapper-cover-image ProductCardVertical__layout-cover-image", 1}, {"a", "class", "ProductCardVertical__image-link  Link js--Link Link_type_default", 1}, {"div", "class", "ProductCardVertical__image-wrapper", 1}, {"div", "class", "ProductCardVertical__picture-container ", 1}, {"img", "class", "ProductCardVertical__picture js--ProductCardInListing__picture", 1}},
		{{"div", "class", "ProductCardVerticalLayout ProductCardVertical__layout", 1}, {"div", "class", "ProductCardVerticalLayout__header", 1}, {"div", "class", "ProductCardVerticalLayout__wrapper-description ProductCardVertical__layout-description", 1}, {"div", "class", "ProductCardVertical__description ", 1}, {"a", "class", " ProductCardVertical__name  Link js--Link Link_type_default", 1}},
		{{"div", "class", "ProductCardVerticalLayout ProductCardVertical__layout", 1}, {"div", "class", "ProductCardVerticalLayout__footer", 1}, {"div", "class", "ProductCardVerticalLayout__wrapper-price", 1}, {"div", "class", "ProductCardVertical_mobile", 1}, {"div", "class", "ProductCardVertical__price-with-amount", 1}, {"div", "class", "ProductPrice ProductPrice_default ProductCardVerticalPrice__price-current", 1}},
		{{"div", "class", "ProductCardVerticalLayout ProductCardVertical__layout", 1}, {"div", "class", "ProductCardVerticalLayout__header", 1}, {"div", "class", "ProductCardVerticalLayout__wrapper-meta", 1}, {"div", "class", "ProductCardVerticalMeta ", 1}, {"div", "class", "ProductCardVerticalMeta__info", 1}, {"div", "class", "Tooltip ProductCardVerticalMeta__tooltip js--Tooltip  js--Tooltip_hover Tooltip_placement_top", 1}, {"a", "class", "  Link js--Link Link_type_default", 1}, {"div", "class", "ProductCardVerticalMeta__text IconWithCount js--IconWithCount", 1}},
		{{"div", "class", "ProductCardVerticalLayout ProductCardVertical__layout", 1}, {"div", "class", "ProductCardVerticalLayout__header", 1}, {"div", "class", "ProductCardVerticalLayout__wrapper-meta", 1}, {"div", "class", "ProductCardVerticalMeta ", 1}, {"div", "class", "ProductCardVerticalMeta__info", 2}, {"div", "class", "Tooltip ProductCardVerticalMeta__tooltip js--Tooltip  js--Tooltip_hover Tooltip_placement_top", 1}, {"a", "class", "  Link js--Link Link_type_default", 1}},
	}
)
