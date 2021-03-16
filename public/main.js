const headerTemplate = Handlebars.templates['components/Header/Header'];
const mainPageTemplate = Handlebars.templates['views/MainPage/MainPage'];

const app = document.getElementById('app');

app.innerHTML += headerTemplate();
app.innerHTML += mainPageTemplate();
