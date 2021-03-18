export class HttpRequest {
    constructor(host = '') {
        this.host = host;
    }

    async get(url = '', query = {}, headers = {}) {
        return await this._fetch({
            url: url,
            query: query,
            headers: headers,
            method: 'GET'
        });
    }

    async post(url = '', body = {}, query = {}, headers = {'Content-Type': 'application/json'}) {
        return await this._fetch({
            url: url,
            query: query,
            headers: headers,
            body: body,
            method: 'POST'
        });
    }

    async delete(url = '', query = {}, headers = {}) {
        return await this._fetch({
            url: url,
            query: query,
            headers: headers,
            method: 'DELETE'
        });
    }

    async _fetch({url, query, body, headers, method}) {
        const encodeQueryString = (query) => {
            if (query === null || typeof query === 'undefined' || Object.keys(query).length === 0) {
                return '';
            }

            let urlParts = [];
            for (const [key, value] of Object.entries(query)) {
                urlParts.push(encodeURIComponent(key) + '=' + encodeURIComponent(value));
            }

            return '?' + urlParts.join('&');
        };

        url += encodeQueryString(query);

        let response = {};

        try {
           response = await fetch(this.host + url, {
                method: method,
                mode: 'no-cors',
                credentials: 'include',
                headers: headers,
                body: body
            });
        } catch (error) {
            return {
                status: 500,
                body: {}
            };
        }

        let parsedBody = {};
        try {
            parsedBody = await response.json();
        } catch (error) {
            parsedBody = {'error': error};
        }

        const parsedResponse = {
            status: response.status,
            body: parsedBody
        };

        return parsedResponse;
    }

}
