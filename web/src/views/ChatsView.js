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

    const link = document.createElement("a")
    link.setAttribute("href", `/chat/${user.id}`)
    link.setAttribute("data-link", "")
    link.innerText = `${user.firstName} ${user.lastName}`

    el.append(link)

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
            <div id="chats"></div>

            Online users:<br>
            <div id="online-users"></div>
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