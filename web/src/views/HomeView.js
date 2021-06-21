import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Home");
    }

    async getHtml() {
        const authorized = Boolean(localStorage.getItem("accessToken"))

        return `
            <nav class="nav">
                <a href="/" class="nav__link" data-link>home</a>
                `
            + (
                authorized ?
                    `<a href="/chats" class="nav__link" data-link>chats</a>`
                    :
                    `<a href="/sign-up" class="nav__link" data-link>sign-up</a>
                    <a href="/sign-in" class="nav__link" data-link>sign-in</a>`
            )
            +
            `
            </nav>
            <h1>Hello</h1>
        `;
    }

}