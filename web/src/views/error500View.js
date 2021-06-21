import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("500");
    }

    async getHtml() {
        return `
            <h1>500 Internal Server Error</h1>
        `;
    }

}