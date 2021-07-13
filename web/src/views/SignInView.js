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

    const data = await fetcher.post(path, body)
    if (data.error) {
        utils.drawErrorMessage(data.error)
        return
    }
    
    localStorage.setItem("accessToken", data.accessToken)
    localStorage.setItem("refreshToken", data.refreshToken)

    const payload = utils.parseJwt(data.accessToken)
    localStorage.setItem("sub", parseInt(payload.sub))
    localStorage.setItem("role", parseInt(payload.role))

    Ws.connect().then(() => {
        router.navigateTo("/")
    })
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
                <input type="password" id="password" placeholder="Password" required>

                <div class="error" id="error-message"></div>

                <button>Sign in</button>
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