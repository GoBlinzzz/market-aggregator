import {HttpRequest} from "../../modules/http.js";

const http = new HttpRequest('/api');
const app = document.getElementById('app');
const submitButton = document.getElementsByClassName('header__search-button').item(0);
const logo = document.getElementsByClassName('header__logo-content').item(0);

const mainPageTemplate = Handlebars.templates['views/MainPage/MainPage'];
const catalogPageTemplate = Handlebars.templates['views/CatalogPage/CatalogPage'];

submitButton.addEventListener('click', sendSearchRequest);
logo.addEventListener('click', goMainPage);

async function sendSearchRequest() {
    const searchString = document.getElementsByClassName('header__search-input').item(0).value;
    if (searchString !== '') {
        http.get({
            'url': '/search' + '?text=' + searchString
        })
            .then((res) => {
                if (res.status === 200) {
                    let catalogCtx = res.body;
                    catalogCtx.request = searchString;

                    const header = document.getElementsByClassName('header').item(0);
                    header.nextElementSibling.remove();

                    const content = catalogPageTemplate(res.body);
                    header.insertAdjacentHTML('afterend', content);
                }
            })
    }
}

function goMainPage() {
    const header = document.getElementsByClassName('header').item(0);
    header.nextElementSibling.remove();

    const content = mainPageTemplate();
    header.insertAdjacentHTML('afterend', content);
}