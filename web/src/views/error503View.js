import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("503");
    }

    async getHtml() {
        return `
            <h1>503 Service Unavailable</h1>
        `;
    }

}