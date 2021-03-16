import {HttpRequest} from "../../modules/http.js";

const http = new HttpRequest('/api');
const app = document.getElementById('app');
const submitButton = document.getElementsByClassName('header__search-button').item(0);

const mainPageTemplate = Handlebars.templates['views/MainPage/MainPage'];
const catalogPageTemplate = Handlebars.templates['views/CatalogPage/CatalogPage'];

submitButton.addEventListener('click', sendSearchRequest);

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

                    app.innerHTML = headerTemplate();
                    app.innerHTML += catalogPageTemplate(res.body);
                }
            })
    }
}