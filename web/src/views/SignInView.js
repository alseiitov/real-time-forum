import AbstractView from "./AbstractView.js";
import router from "../index.js"

const signIn = async (username, password) => {
    const url = "http://localhost:8081/api/users/sign-in"

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

        router.navigateTo("/chats")
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
        const signInForm = document.querySelector("#sign-in-form")
        signInForm.addEventListener("submit", function () {
            const username = document.querySelector("#username").value
            const password = document.querySelector("#password").value

            signIn(username, password)
        })
    }
}