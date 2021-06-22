import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("401");
    }

    async getHtml() {
        return `
            <h1>401 Unauthorized</h1>
        `;
    }

}