import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params, user) {
        super(params);
        this.user = user
    }

    async getHtml() {
        const authorized = Boolean(this.user.id)

        return `
            <a href="/" class="nav__link" data-link>home</a>
            `
            +
            (authorized ?
                `
                <a href="/chats" class="nav__link" data-link>chats</a>
                `
                :
                `
                <a href="/sign-up" class="nav__link" data-link>sign-up</a>
                <a href="/sign-in" class="nav__link" data-link>sign-in</a>
                `
            )

    }

}