import {HttpRequest} from "../../modules/http.js";

window.http = new HttpRequest('/api');
window.searchString = '';

const app = document.getElementById('app');
const submitButton = document.getElementsByClassName('header__search-button')
    .item(0);
const logo = document.getElementsByClassName('header__logo-content').item(0);
const searchInput = document.getElementsByClassName('header__search-input')
    .item(0);
const shoppingCartButton = document.getElementsByClassName('header__shop-cart-div').item(0);

const mainPageTemplate = Handlebars.templates['views/MainPage/MainPage'];
const catalogPageTemplate = Handlebars.templates['views/CatalogPage/CatalogPage'];
const shoppingCartPageTemplate = Handlebars.templates['views/ShoppingCartPage/ShoppingCartPage'];

window.searchRequestPending = false;

submitButton.addEventListener('click', sendSearchRequest);

logo.addEventListener('click', goMainPage);

searchInput.addEventListener('keyup', (event) => {
    if (event.keyCode === 13) {
        submitButton.click();
    }
});

shoppingCartButton.addEventListener('click', () => {
    const ctx = {
        items: [
            {
                image: 'https://items.s1.citilink.ru/1380949_v01_b.jpg',
                title: 'Смартфон XIAOMI Redmi Note 9 Pro 6/128Gb, серый',
                sourceMarket: 'ctl',
                marketTitle: 'Citilink',
                price: '44000$'
            }
        ]
    };
    goShoppingCartPage(ctx);
});

async function sendSearchRequest(how = '') {
    searchString = document.getElementsByClassName('header__search-input')
        .item(0).value;
    if (!searchRequestPending && searchString !== '') {
        searchRequestPending = true;
        http.get({
            'url': '/search' + '?text=' + searchString + '&how=' + how
        })
            .then((res) => {
                searchRequestPending = false;
                if (res.status === 200) {
                    let catalogCtx = res.body;
                    catalogCtx.items.forEach(item => {
                        if (item.rating === null) {
                            item.rating = {
                                count: 0,
                                withHalf: false
                            }
                        }
                        item.rating.withHalf = Number(item.rating.withHalf);
                        switch (item.sourceMarket) {
                            case 'wb':
                                item.marketTitle = 'Wildberries';
                                item.marketLink = 'https://www.wildberries.ru/';
                                item.reviewLink = item.link;
                                break;
                            case 'ctl':
                                item.marketTitle = 'Citilink';
                                item.marketLink = 'https://citilink.ru/';
                                item.reviewLink = item.link + 'otzyvy';
                                break;
                        }
                    })
                    catalogCtx.request = searchString;

                    goCatalogPage(catalogCtx);
                }
            });
    }
}

function goMainPage() {
    const header = document.getElementsByClassName('header').item(0);
    header.nextElementSibling.remove();

    const content = mainPageTemplate();
    header.insertAdjacentHTML('afterend', content);
}

function goCatalogPage(ctx) {
    const header = document.getElementsByClassName('header')
        .item(0);
    header.nextElementSibling.remove();

    const content = catalogPageTemplate(ctx);
    header.insertAdjacentHTML('afterend', content);
}

function goShoppingCartPage(ctx) {
    const header = document.getElementsByClassName('header')
        .item(0);
    header.nextElementSibling.remove();

    const content = shoppingCartPageTemplate(ctx);
    header.insertAdjacentHTML('afterend', content);
}

window.sendSearchRequest = sendSearchRequest;
