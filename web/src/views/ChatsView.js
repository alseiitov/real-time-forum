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


const newUserLinkElement = (user) => {
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

    const avatatEl = document.createElement("div")
    avatatEl.innerHTML = `<img src="http://${API_HOST_NAME}/images/${chat.user.avatar}">`;
    el.append(avatatEl)

    const messageEl = document.createElement("div")

    const link = newUserLinkElement(chat.user)
    messageEl.append(link)

    if (chat.lastMessage) {
        const lastMessage = document.createElement("p")
        lastMessage.innerText = `${chat.lastMessage.message}`
        messageEl.append(lastMessage)

        const lastMessageDate = document.createElement("p")
        lastMessageDate.innerText = `${new Date(chat.lastMessage.date).toLocaleString()}`
        messageEl.append(lastMessageDate)
    }

    el.append(messageEl)

    return el
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Chats");
    }

    async getHtml() {
        return `
            <h3>Your chats:</h3>
            <div id="chats"></div>

            <h3>Online users:</h3>
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
                    const el = newChatElement({user: user})
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
                console.log(chat)
            })
        }
    }
}