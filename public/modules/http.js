export class HttpRequest {
    constructor(host = '') {
        this.host = host;
    }

    async get(params = {url: ''}) {
        return await this.#fetchSafe(params);
    }

    async post(params = {
        url: '', data: {}, headers: {}
    }) {
        return await this.#fetchUnsafe(params);
    }

    async put(params = {url: '', data: {}}) {
        return await this.#fetchUnsafe({...params, method: 'PUT'})
    }

    async delete(params = {url: ''}) {
        return await this.#fetchUnsafe({...params, method: 'DELETE'})
    }

    async #fetchSafe({url = '', method = 'GET'}) {
        const response = await fetch(this.host + url, {
            method: method,
            mode: 'no-cors',
            credentials: 'include'
        });

        let parsedResponse = {status: response.status};
        let parsedBody = null;

        try {
            parsedBody = await response.json();
        } catch (error) {
            parsedBody = {"error": error};
        }

        parsedResponse.body = parsedBody;

        return parsedResponse;
    }

    async #fetchUnsafe({
                           url = '', data = JSON.stringify({}), headers = {
            'Content-Type': 'application/json'
        }, method = 'POST'
                       }) {
        const response = await fetch(this.host + url, {
            method: method,
            mode: 'no-cors',
            credentials: 'include',
            headers: headers,
            body: data
        });

        let parsedResponse = {status: response.status};
        let parsedBody = null;

        try {
            parsedBody = await response.json();
        } catch (error) {
            parsedBody = {"error": error};
        }

        parsedResponse.body = parsedBody;

        return parsedResponse;
    }

}