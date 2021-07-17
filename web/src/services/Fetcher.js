import Utils from "./Utils.js"
import Router from "./../index.js"

const fetcher = {
    get: async (path, body) => {
        return makeRequest(path, body, "GET")
    },

    post: async (path, body) => {
        return makeRequest(path, body, "POST")
    },

    refreshToken: async () => {
        const accessToken = localStorage.getItem("accessToken")
        const refreshToken = localStorage.getItem("refreshToken")

       localStorage.removeItem("accessToken")
       localStorage.removeItem("refreshToken")
        
        const path = "/api/auth/refresh"
        const body = { accessToken: accessToken, refreshToken: refreshToken }
        const data = await fetcher.post(path, body)
        const payload = Utils.parseJwt(data.accessToken)
        
        localStorage.setItem("accessToken", data.accessToken)
        localStorage.setItem("refreshToken", data.refreshToken)
        localStorage.setItem("role", parseInt(payload.role))
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
    
    if (response.status == 400) {
        if (respBody.error == "invalid token") {
            Utils.logOut()
            Router.navigateTo("/sign-in")
            return
        }

        const errorEl = document.getElementById('error-message')
        if (errorEl) {
            errorEl.innerText = respBody.error
        } else {
            Utils.showError(response.status, respBody.error)
        }
        return
    }

    if (response.status == 401) {
        if (respBody.error == "token has expired") {
            await fetcher.refreshToken()
            return await makeRequest(path, body, method)
        }
        return respBody
    }

    if (response.status == 409) {
        return respBody
    }

    if (!response.ok) {
        Utils.showError(response.status, respBody.error)
        return
    }

    return respBody
}

export default fetcher