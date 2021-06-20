import AbstractView from "./AbstractView.js";
import router from "../index.js"
import utils from "../services/Utils.js"

const signUp = async (input) => {
    const url = `http://${API_HOST_NAME}/api/users/sign-up`

    const options = {
        method: "POST",
        body: JSON.stringify(input)
    }


    const response = await fetch(url, options)
    // if (response.ok) {
    //     const data = await response.json()

    //     localStorage.setItem("accessToken", data.accessToken)
    //     localStorage.setItem("refreshToken", data.refreshToken)

    //     const payload = utils.parseJwt(data.accessToken)
    //     localStorage.setItem("sub", parseInt(payload.sub))
    //     localStorage.setItem("role", parseInt(payload.role))

    //     router.navigateTo("/chats")
    // }
    if (response.status != 201) {
        const data = await response.json()
        console.log(data)
    } else {
        router.navigateTo("/sign-in")
    }
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Sign up");
    }

    async getHtml() {
        return `
            <form id="sign-up-form" onsubmit="return false;">
                Username: <br>
                <input type="text" id="username" placeholder="Username" required minlength="2" maxlength="64"> <br> <br>

                First name: <br>
                <input type="text" id="first-name" placeholder="First name" required minlength="2" maxlength="64"> <br> <br>

                Last name: <br>
                <input type="text" id="last-name" placeholder="Last name" required minlength="2" maxlength="64"> <br> <br>

                Age: <br>
                <input type="number" id="age" placeholder="Age" required min="12" max="110"> <br> <br>

                Gender: <br>
                <input type="radio" name="gender" id="gender-male" value="1" required> Male
                <input type="radio" name="gender" id="gender-female" value="2"> Female <br> <br>

                E-mail: <br>
                <input type="email" id="email" placeholder="E-mail" required maxlength="64"> <br><br>
                
                Password: <br>
                <input type="password" id="password" placeholder="Password" maxlength="64" required> <br> <br>
                
                Confirm password: <br>
                <input type="password" id="password-confirm" placeholder="Password" maxlength="64" required> <br> <br>

                <div id="input-error"></div>
                
                <button type="submit">Sign up</button>
            </form>
        `;
    }

    async init() {
        const signUpForm = document.getElementById("sign-up-form")
        const inputError = document.getElementById("input-error")

        signUpForm.addEventListener("submit", function () {
            inputError.innerText = ""

            const password = document.getElementById("password")
            const passwordConfirm = document.getElementById("password-confirm")

            if (password.value != passwordConfirm.value) {
                inputError.innerText = "Passwords Don't Match"
            } else {
                let input = {
                    username: document.getElementById("username").value,
                    firstName: document.getElementById("first-name").value,
                    lastName: document.getElementById("last-name").value,
                    age: parseInt(document.getElementById("age").value),
                    gender: parseInt(document.querySelector('input[name="gender"]:checked').value),
                    email: document.getElementById("email").value,
                    password: password.value,
                }

                signUp(input)
            }
        })
    }
}