import AbstractView from "./AbstractView.js";
import Utils from "../services/Utils.js"
import Router from "../index.js"

export default class extends AbstractView {
    constructor(params, user) {
        super(params);
        this.user = user
    }

    async getHtml() {
        const authorized = Boolean(this.user.id)

        return `
            <a href="/" class="nav__link" id="home-button" data-link>home</a>
            `
            +
            (authorized ?
                `
                <a href="/chats" class="nav__link" id="chats-button" data-link>chats</a>
                <a href="" class="nav__link"id="sign-out-button" data-link>sign-out</a>
                `
                :
                `
                <a href="/sign-up" class="nav__link" id="sign-in-button" data-link>sign-up</a>
                <a href="/sign-in" class="nav__link" id="sign-up-button" data-link>sign-in</a>
                `
            )
    }

    async init() {
        const signOutButton = document.getElementById("sign-out-button")
        signOutButton.addEventListener('click', () => {
            Utils.logOut()
            Router.navigateTo("/")
        })
    }
}