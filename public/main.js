const headerTemplate = Handlebars.templates['components/Header/Header'];
const mainPageTemplate = Handlebars.templates['views/MainPage/MainPage'];
const catalogPageTemplate = Handlebars.templates['views/CatalogPage/CatalogPage'];


const app = document.getElementById('app');

app.innerHTML += headerTemplate();
app.innerHTML += mainPageTemplate();
