import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Home");
    }

    async getHtml() {
        return `
            <h1>Hello</h1>
           
            <p>
                <a href="/chats" data-link>Chats</a>
            </p>
        `;
    }

}