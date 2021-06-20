import Ws from "../services/Ws.js";
import AbstractView from "./AbstractView.js";

var requestOnlineUsersInterval 

const requestOnlineUsers = () => {
    Ws.send(JSON.stringify({ type: "onlineUsersRequest" }))
}

const newUserElement = (user) => {
    const el = document.createElement("div");
    el.classList.add("user")
    el.id = `user-${user.id}`
    el.innerText = `${user.firstName} ${user.lastName}\n`

    return el
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Chats");
    }

    async getHtml() {
        return `
            Your chats:<br>
            <div id="chats">
                <a href="/chat/1" data-link>Chat with user 1</a>
                <a href="/chat/2" data-link>Chat with user 2</a>
            </div>
                Online users:<br>
            <div id="online-users">
            </div>
        `;
    }

    async init() {
        requestOnlineUsers()
        clearInterval(requestOnlineUsersInterval)
        requestOnlineUsersInterval = setInterval(requestOnlineUsers, 10_000)
    }

    static drawOnlineUsers(users) {
        if (users != null) {
            const onlineUsersEl = document.getElementById("online-users");
            if (onlineUsersEl) {
                onlineUsersEl.innerText = ""
                users.forEach((user) => {
                    const el = newUserElement(user)
                    onlineUsersEl.prepend(el)
                })
            }
        }
    }
}