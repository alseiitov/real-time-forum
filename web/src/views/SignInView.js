import AbstractView from "./AbstractView.js";
import router from "../index.js"
import utils from "../services/Utils.js"
import Ws from "../services/Ws.js";
import fetcher from "../services/Fetcher.js"

const path = `/api/users/sign-in`

const signIn = async (usernameOrEmail, password) => {

    let body = {
        usernameOrEmail: usernameOrEmail,
        password: password
    }

    // fetch(url, options).then((response) => {
    //     switch (response.status) {
    //         case 200:
    //             response.json().then((data) => {
    //                 localStorage.setItem("accessToken", data.accessToken)
    //                 localStorage.setItem("refreshToken", data.refreshToken)

    //                 const payload = utils.parseJwt(data.accessToken)
    //                 localStorage.setItem("sub", parseInt(payload.sub))
    //                 localStorage.setItem("role", parseInt(payload.role))

    //                 Ws.connect().then(() => {
    //                     router.navigateTo("/")
    //                 })
    //             })
    //             break
    //         case 400: case 401:
    //             response.json().then((data) => {
    //                 drawError(data.error)
    //             })
    //             break
    //         case 500:
    //             router.navigateTo("/500")
    //             break
    //     }
    // })

    const response = await fetcher.post(path, body)

    switch (response.status) {
        case 200:
            response.json().then((data) => {
                localStorage.setItem("accessToken", data.accessToken)
                localStorage.setItem("refreshToken", data.refreshToken)

                const payload = utils.parseJwt(data.accessToken)
                localStorage.setItem("sub", parseInt(payload.sub))
                localStorage.setItem("role", parseInt(payload.role))

                Ws.connect().then(() => {
                    router.navigateTo("/")
                })
            })
            break
        case 400: case 401:
            response.json().then((data) => {
                drawError(data.error)
            })
            break
    }
}

const drawError = (err) => {
    const inputError = document.getElementById("input-error")
    inputError.innerText = err
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Sign in");
    }

    async getHtml() {
        return `
            <form id="sign-in-form" onsubmit="return false;">
                Username or email: <br>
                <input type="text" id="usernameOrEmail" placeholder="Username or email" required> <br> <br>

                Password: <br>
                <input type="password" id="password" placeholder="Password" required> <br> <br>

                <div id="input-error"></div>

                <button class="button">Sign in</button>
            </form>

        `;
    }

    async init() {
        const signInForm = document.getElementById("sign-in-form")
        signInForm.addEventListener("submit", function () {
            const usernameOrEmail = document.getElementById("usernameOrEmail").value
            const password = document.getElementById("password").value

            signIn(usernameOrEmail, password)
        })
    }
}