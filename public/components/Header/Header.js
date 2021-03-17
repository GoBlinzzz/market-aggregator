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
    let ctx = {};
    sendCartRequest().then((res) => {
        if (res.status == 200) {
            ctx = res.body;
            ctx.items.forEach(item => {
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
                    case 'eld':
                        item.marketTitle = 'Eldorado';
                        item.marketLink = 'https://www.eldorado.ru/';
                        item.reviewLink = item.link + '?show=response#customTabAnchor';
                }
            });
        }
        goShoppingCartPage(ctx);
    });
});

async function sendSearchRequest(how = '') {
    searchString = document.getElementsByClassName('header__search-input')
        .item(0).value;
    if (!window.searchRequestPending && searchString !== '') {
        window.searchRequestPending = true;
        const query = {
            'text': searchString,
            'how': how
        }
        await http.get('/search', query)
            .then((res) => {
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
                            case 'eld':
                                item.marketTitle = 'Eldorado';
                                item.marketLink = 'https://www.eldorado.ru/';
                                item.reviewLink = item.link + '?show=response#customTabAnchor';
                        }
                    });
                    catalogCtx.request = searchString;

                    goCatalogPage(catalogCtx);
                }
            });
        window.searchRequestPending = false;
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

async function sendAddToCartRequest(id) {
    const elem = document.getElementById(id).parentElement.parentElement.parentElement;
    const item = {
        link: id,
        sourceMarket: elem.getElementsByClassName('catalog__item-source-market').item(0).textContent,
        image: elem.getElementsByClassName('catalog__item-img').item(0).getAttribute('src'),
        title: elem.getElementsByClassName('catalog__item-title').item(0).textContent,
        price: elem.getElementsByClassName('catalog__item-price').item(0).textContent
    };

    switch (item.sourceMarket) {
        case 'Citilink':
            item.sourceMarket = 'ctl';
            break;
        case 'Eldorado':
            item.sourceMarket = 'eld';
            break;
        case 'Wildberries':
            item.sourceMarket = 'wb';
            break;
    }

    http.post('/add-to-cart', JSON.stringify(item));
}

async function sendCartRequest() {
    return await http.get('/cart');
}

function deleteFromCart(id) {
    const elem = document.getElementById(id);
    const itemPos = Array.from(document.getElementsByClassName('shopping-cart__item-div')).indexOf(elem);

    sendDeleteFromCartRequest(itemPos).then((res) => {
        if (res.status === 200) {
            shoppingCartButton.click();
        }
    });
}

async function sendDeleteFromCartRequest(itemPos) {
    return await http.delete('/cart', {'id': itemPos});
}

window.sendSearchRequest = sendSearchRequest;
window.sendAddToCartRequest = sendAddToCartRequest;
window.deleteFromCart = deleteFromCart;