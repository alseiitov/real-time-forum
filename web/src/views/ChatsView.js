import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Chats");
    }

    async getHtml() {
        return `
            <a href="/chat/1" data-link>Chat with user 1</a>
            <a href="/chat/2" data-link>Chat with user 2</a>
        `;
    }
}