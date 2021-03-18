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
const header = document.getElementsByClassName('header')
    .item(0);

const mainPageTemplate = Handlebars.templates['views/MainPage/MainPage'];
const catalogPageTemplate = Handlebars.templates['views/CatalogPage/CatalogPage'];
const shoppingCartPageTemplate = Handlebars.templates['views/ShoppingCartPage/ShoppingCartPage'];
const waitingTemplate = Handlebars.templates['views/Waiting/Waiting'];

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
        startWaiting();
        if (res.status == 200) {
            ctx = res.body;
            if (ctx.items !== null && ctx.items.length !== 0) {
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
            } else {
                ctx.empty = true;
            }
        }
        goShoppingCartPage(ctx);
    });
});

async function sendSearchRequest(how = '') {
    searchString = document.getElementsByClassName('header__search-input')
        .item(0).value;
    if (!window.searchRequestPending && searchString !== '') {
        window.searchRequestPending = true;
        startWaiting();

        const query = {
            'text': searchString,
            'how': how
        }

        let catalogCtx = {};

        console.log(how)

        await http.get('/search', query)
            .then((res) => {
                if (res.status === 200) {
                    catalogCtx = res.body;
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

                } else if (res.status === 404) {
                    catalogCtx = {count: 0, notFound: true};
                }

                catalogCtx.request = searchString;
                switch (how) {
                    case 'rating':
                        catalogCtx.rating = true;
                        break;
                    case 'aprice':
                        catalogCtx.aprice = true;
                        break;
                    case 'dprice':
                        catalogCtx.dprice = true;
                        break;
                    default:
                        catalogCtx.popularity = true;
                        break;
                }
                goCatalogPage(catalogCtx);
            });
        window.searchRequestPending = false;
    }
}

function goMainPage() {
    header.nextElementSibling.remove();

    const content = mainPageTemplate();
    header.insertAdjacentHTML('afterend', content);
}

function goCatalogPage(ctx) {
    header.nextElementSibling.remove();

    const content = catalogPageTemplate(ctx);
    header.insertAdjacentHTML('afterend', content);
}

function goShoppingCartPage(ctx) {
    header.nextElementSibling.remove();

    const content = shoppingCartPageTemplate(ctx);
    header.insertAdjacentHTML('afterend', content);
}

function startWaiting() {
    header.nextElementSibling.remove();

    const content = waitingTemplate();
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
    return await http.post('/delete-from-cart', {},{'id': itemPos});
}

window.sendSearchRequest = sendSearchRequest;
window.sendAddToCartRequest = sendAddToCartRequest;
window.deleteFromCart = deleteFromCart;