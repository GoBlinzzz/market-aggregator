(function () {
    function goLink(link) {
        Object.assign(document.createElement('a'), {
            target: '_blank',
            href: link,
        }).click();
    }

    window.goLink = goLink;
})()