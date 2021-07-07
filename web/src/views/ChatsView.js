import intervals from "../services/Intervals.js";
import Ws from "../services/Ws.js";
import AbstractView from "./AbstractView.js";

var requestOnlineUsersInterval

const requestOnlineUsers = () => {
    Ws.send(JSON.stringify({ type: "onlineUsersRequest" }))
}

const requestChats = () => {
    Ws.send(JSON.stringify({ type: "chatsRequest" }))
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


const newChatElement = (chat) => {
    const el = document.createElement("div");
    el.classList.add("chat")
    el.id = `chat-${chat.user.id}`

    const link = newUserElement(chat.user)
    const lastMessage = document.createElement("p")
    lastMessage.innerText = `${chat.lastMessage.message}`

    const lastMessageDate = document.createElement("p")
    lastMessageDate.innerText = `${new Date(chat.lastMessage.date).toLocaleString()}`

    el.append(link)
    el.append(lastMessage)
    el.append(lastMessageDate)

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
        requestChats()
        requestOnlineUsers()

        intervals.clear(requestOnlineUsersInterval)
        requestOnlineUsersInterval = setInterval(requestOnlineUsers, 10_000)
        intervals.add(requestOnlineUsersInterval)
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

    static drawChats(chats) {
        if (chats != null) {
            const chatsEl = document.getElementById("chats");
            chatsEl.innerHTML = ""
            chats.forEach((chat) => {
                const el = newChatElement(chat)
                chatsEl.append(el)
            })
        }
    }
}