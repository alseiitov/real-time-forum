import Utils from "./Utils.js"

const fetcher = {
    get: async (path, body) => {
        return makeRequest(path, body, "GET")
    },

    post: async (path, body) => {
        return makeRequest(path, body, "POST")
    }
}

const makeRequest = async (path, body, method) => {
    const url = `http://${API_HOST_NAME}${path}`

    const options = {
        mode: 'cors',
        method: method,
        body: JSON.stringify(body)
    }

    const accessToken = localStorage.getItem("accessToken")
    if (accessToken != undefined) {
        options.headers = new Headers({
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${accessToken}`,
        })
    }

    const response = await fetch(url, options).catch((e) => {
        Utils.showError(503)
        return
    })

    var respBody

    try {
        respBody = await response.json()
    } catch {
        return
    }

    if (response.status == 401 || response.status == 409) {
        return respBody
    }

    if (!response.ok) {
        Utils.showError(response.status, respBody.error)
        return
    }

    return respBody
}

export default fetcher