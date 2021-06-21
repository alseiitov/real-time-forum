import router from "../index.js"

const fetcher = {
    get: async (path, body) => {
        return await makeRequest(path, body, "GET")
    },

    post: async (path, body) => {
        return await makeRequest(path, body, "POST")
    }
}

const makeRequest = async (path, body, method) => {
    const url = `http://${API_HOST_NAME}${path}`

    const options = {
        method: method,
        body: JSON.stringify(body)
    }

    const response = await fetch(url, options).catch(() => {
        router.navigateTo("/503")
        return
    })

    if (response.status == 404) {
        router.navigateTo("/404")
        return
    }

    if (response.status == 500) {
        router.navigateTo("/500")
        return
    }

    return response
}

export default fetcher