import AbstractView from "./AbstractView.js";
import Utils from "../services/Utils.js"
import Ws from "../services/Ws.js";
import intervals from "../services/Intervals.js";


export default class extends AbstractView {
    constructor(params, user) {
        super(params);
        this.user = user
    }

    async getHtml() {
        const authorized = Boolean(this.user.id)

        return `
            <a href="/" class="nav__link" id="home-button" data-link>Home</a>
            `
            +
            (authorized ?
                `
                <a href="/new-post" class="nav__link" id="new-post-button" data-link>New post</a>
                <a href="/chats" class="nav__link" id="chats-button" data-link>Chats</a>
                <a href="/user/${this.user.id}" class="nav__link" id="chats-button" data-link>Profile</a>
                <a href="/" class="nav__link" id="sign-out-button" data-link>Logout</a>
                `
                :
                `
                <a href="/sign-up" class="nav__link" id="sign-in-button" data-link>Sign up</a>
                <a href="/sign-in" class="nav__link" id="sign-up-button" data-link>Sign in</a>
                `
            )
    }

    async init() {
        const signOutButton = document.getElementById("sign-out-button")
        if (signOutButton) {
            signOutButton.addEventListener('click', () => {
                Utils.logOut()
                Ws.disconnect()
                intervals.clearAll()
            })
        }
    }
}