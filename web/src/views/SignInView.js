import AbstractView from "./AbstractView.js";
import router from "../index.js"
import utils from "../services/Utils.js"

const signIn = async (username, password) => {
    const url = `http://${API_HOST_NAME}/api/users/sign-in`

    const options = {
        method: "POST",
        body: JSON.stringify({
            usernameOrEmail: username,
            password: password
        })
    }

    const response = await fetch(url, options)
    if (response.ok) {
        const data = await response.json()

        localStorage.setItem("accessToken", data.accessToken)
        localStorage.setItem("refreshToken", data.refreshToken)

        const payload = utils.parseJwt(data.accessToken)
        localStorage.setItem("sub", parseInt(payload.sub))
        localStorage.setItem("role", parseInt(payload.role))

        router.navigateTo("/")
    }
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
                <input type="text" id="username" placeholder="Username or email" required> <br> <br>
                Password: <br>
                <input type="password" id="password" placeholder="Password" required> <br> <br>
                <button class="button">Sign in</button>
            </form>
        `;
    }

    async init() {
        const signInForm = document.getElementById("sign-in-form")
        signInForm.addEventListener("submit", function () {
            const username = document.getElementById("username").value
            const password = document.getElementById("password").value

            signIn(username, password)
        })
    }
}