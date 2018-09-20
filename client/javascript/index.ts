import fetch from 'cross-fetch';

const parseResponse = (value: string) : string|boolean => {
    if (value === "on") {
        return true;
    }
    return value;
}

class Client {
    server: string
    default: boolean|string

    constructor(server: string, def: boolean|string = true) {
        this.server = server;
        this.default = def;
    }

    async get(flag: string, environment: string): Promise<string|boolean> {
        return await fetch(`${this.server}/flags/${flag}/${environment}`).then(response => {
            if (response.ok) {
                return response.text();
            }
            throw new Error("Bad response");
        }).then(response => parseResponse(response))
        .catch(err => {
            return this.default;
        })
    }
}

export default Client;