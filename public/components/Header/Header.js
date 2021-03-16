import {HttpRequest} from "../../modules/http.js";

window.http = new HttpRequest('/api');
window.searchString = '';

const app = document.getElementById('app');
const submitButton = document.getElementsByClassName('header__search-button')
    .item(0);
const logo = document.getElementsByClassName('header__logo-content').item(0);
const searchInput = document.getElementsByClassName('header__search-input')
    .item(0);

const mainPageTemplate = Handlebars.templates['views/MainPage/MainPage'];
const catalogPageTemplate = Handlebars.templates['views/CatalogPage/CatalogPage'];

window.searchRequestPending = false;

submitButton.addEventListener('click', sendSearchRequest);
logo.addEventListener('click', goMainPage);
searchInput.addEventListener('keyup', (event) => {
    if (event.keyCode === 13) {
        submitButton.click();
    }
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
                        if (!item.rating) {
                            item.rating = '0';
                        }
                        if (!item.reviewCount) {
                            item.reviewCount = '0';
                        }
                        switch (item.sourceMarket) {
                            case 'wb':
                                item.marketTitle = 'Wildberries';
                                item.marketLink = 'https://www.wildberries.ru/';
                                break;
                            case 'ctl':
                                item.marketTitle = 'Citilink';
                                item.marketLink = 'https://citilink.ru/';
                                break;
                        }
                    })
                    catalogCtx.request = searchString;

                    const header = document.getElementsByClassName('header')
                        .item(0);
                    header.nextElementSibling.remove();

                    const content = catalogPageTemplate(res.body);
                    header.insertAdjacentHTML('afterend', content);
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

window.sendSearchRequest = sendSearchRequest;
