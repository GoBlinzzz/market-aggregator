const headerTemplate = Handlebars.templates['components/Header/Header'];
const mainPageTemplate = Handlebars.templates['views/MainPage/MainPage'];
const catalogPageTemplate = Handlebars.templates['views/CatalogPage/CatalogPage'];


const app = document.getElementById('app');

const catalogCtx = {
    'count': 16,
    'items': [
        {
            'image': 'https://images.wbstatic.net/c252x336/new/8010000/8018536-1.jpg',
            'price': '1450$',
            'title': 'Gigabyte / Видеокарта nVidia GeForce GT 730/2048 Мб/GDDR5',
            'reviewCount': '2'
        },
        {
            'image': 'https://images.wbstatic.net/c252x336/new/8010000/8018536-1.jpg',
            'price': '1450$',
            'title': 'Gigabyte / Видеокарта nVidia GeForce GT 730/2048 Мб/GDDR5',
            'reviewCount': '2'
        },
        {
            'image': 'https://images.wbstatic.net/c252x336/new/8010000/8018536-1.jpg',
            'price': '1450$',
            'title': 'Gigabyte / Видеокарта nVidia GeForce GT 730/2048 Мб/GDDR5',
            'reviewCount': '2'
        },
        {
            'image': 'https://images.wbstatic.net/c252x336/new/8010000/8018536-1.jpg',
            'price': '1450$',
            'title': 'Gigabyte / Видеокарта nVidia GeForce GT 730/2048 Мб/GDDR5',
            'reviewCount': '2'
        },
        {
            'image': 'https://images.wbstatic.net/c252x336/new/8010000/8018536-1.jpg',
            'price': '1450$',
            'title': 'Gigabyte / Видеокарта nVidia GeForce GT 730/2048 Мб/GDDR5',
            'reviewCount': '2'
        },
    ],

};

app.innerHTML += headerTemplate();
// app.innerHTML += mainPageTemplate;
app.innerHTML += catalogPageTemplate(catalogCtx);