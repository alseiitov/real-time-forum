import intervals from "../services/Intervals.js";
import Ws from "../services/Ws.js";
import AbstractView from "./AbstractView.js";
import Utils from "../services/Utils.js"


var requestOnlineUsersInterval
var recipientID


var loadMessages

const requestOnlineUsers = () => {
    Ws.send(JSON.stringify({ type: "onlineUsersRequest" }))
}

const requestChats = () => {
    Ws.send(JSON.stringify({ type: "chatsRequest" }))
}

const newMessageElement = (message) => {
    const el = document.createElement("div");
    el.id = `message-${message.id}`

    const userID = parseInt(localStorage.getItem("sub"))
    if (!message.read && message.senderID == recipientID) {
        Ws.send(JSON.stringify({ type: "readMessageRequest", body: { messageID: message.id } }))
        changeChatUnreadCount(document.getElementById(`chat-${message.senderID}-unread-messages-count`), -1)
    }

    const messageText = document.createElement('p')
    messageText.classList.add('message-text')
    messageText.innerText = message.message

    const messageInfo = document.createElement('div')
    messageInfo.classList.add('message-info')
    messageInfo.id = `message-${message.id}-info`

    const messageDate = document.createElement('p')
    messageDate.classList.add('message-date')
    messageDate.innerText = new Date(Date.parse(message.date)).toLocaleString()

    const readStatus = document.createElement('p')
    readStatus.classList.add('message-status')
    readStatus.id = `message-${message.id}-status`

    if (!message.read && message.senderID == userID) {
        readStatus.innerText = '✓'
    } else {
        readStatus.innerText = '✓✓'
    }

    el.classList.add("message")
    if (message.senderID == userID) {
        el.classList.add("sended-message")
    } else {
        el.classList.add("received-message")
    }

    messageInfo.append(messageDate)
    messageInfo.append(readStatus)

    el.append(messageText)
    el.append(messageInfo)

    return el
}


const newChatElement = (chat) => {
    const el = document.createElement("div");
    el.classList.add("chat")
    el.id = `chat-${chat.user.id}`
    el.style.cursor = 'pointer'

    if (chat.user.id == recipientID) {
        el.classList.add('active')
    }

    const avatatEl = document.createElement("div")
    avatatEl.innerHTML = `<img src="http://${API_HOST_NAME}/images/${chat.user.avatar}">`;
    el.append(avatatEl)

    const messageEl = document.createElement("div")
    messageEl.classList.add('chat-info')

    const name = document.createElement("p")
    name.innerText = `${chat.user.firstName} ${chat.user.lastName}`
    messageEl.append(name)

    const lastMessage = document.createElement("p")
    lastMessage.id = `chat-${chat.user.id}-lastMessage`

    const lastMessageDate = document.createElement("p")
    lastMessageDate.id = `chat-${chat.user.id}-lastMessageDate`

    if (chat.lastMessage) {
        lastMessage.innerText = `${chat.lastMessage.message}`
        lastMessageDate.innerText = `${new Date(chat.lastMessage.date).toLocaleString()}`
    }

    messageEl.append(lastMessage)
    messageEl.append(lastMessageDate)

    const unreadMessagesCount = document.createElement("div")
    unreadMessagesCount.classList.add(`chat-unread-messages-count`)
    unreadMessagesCount.id = `chat-${chat.user.id}-unread-messages-count`
    changeChatUnreadCount(unreadMessagesCount, chat.unreadMessagesCount)

    el.append(messageEl)
    el.append(unreadMessagesCount)

    el.addEventListener("click", () => {
        Array.from(document.getElementsByClassName("chat")).forEach(el => { el.classList.remove('active') })
        Array.from(document.getElementsByClassName("chat-unread-messages-count")).forEach(el => { el.classList.remove('active') })
        
        el.classList.add('active')
        unreadMessagesCount.classList.add('active')

        document.getElementById("message-form").style.display = "block"
        document.getElementById("message-input").value = ""
        document.getElementById("chat-messages").innerHTML = ""
        recipientID = chat.user.id
        Ws.send(JSON.stringify({ type: "messagesRequest", body: { userID: recipientID, lastMessageID: 0 } }))
    })

    return el
}

const changeChatUnreadCount = (el, n) => {
    n += parseInt(el.innerText) || 0
    el.innerText = n

    if (n > 0) {
        el.style.opacity = 100
    } else {
        el.style.opacity = 0
    }
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Chats");
    }

    async getHtml() {
        return `
            <div id="chats-container">
                <div id="chats">
                    <p class="chats-title">Your chats:</p>              
                    <div id="users-chats"></div>

                    <p class="chats-title">Online users:</p>
                    <div id="online-users"></div>
                </div>
                <div>
                    <div id="chat-messages"></div>
                    <form id="message-form">
                        <input type="text" id="message-input" size="64" placeholder="Send message" autocomplete="off" autofocus/>
                    </form>
                </div>
            </div>
        `;
    }

    async init() {
        const chatMessages = document.getElementById("chat-messages");
        const messageForm = document.getElementById("message-form");
        messageForm.style.display = 'none'

        const messageInput = document.getElementById("message-input");

        loadMessages = Utils.debounce(function () {
            if (chatMessages.scrollTop < chatMessages.scrollHeight * 0.1) {
                let offsetMsg = document.querySelector('.message')
                let offsetMsgID = parseInt(offsetMsg.id.split('-')[1]);
                Ws.send(JSON.stringify({ type: "messagesRequest", body: { userID: recipientID, lastMessageID: offsetMsgID } }))
            }

            if (chatMessages.scrollTop == 0) {
                chatMessages.scrollTop = 1
            }
        }, 100, true)

        chatMessages.addEventListener("scroll", loadMessages)

        messageForm.onsubmit = function () {
            if (!messageInput.value) {
                return false;
            }

            let message = { recipientID: recipientID, message: messageInput.value }
            Ws.send(JSON.stringify({ type: "message", body: message }));
            messageInput.value = "";
            return false;
        };

        requestChats()
        requestOnlineUsers()

        intervals.clear(requestOnlineUsersInterval)
        requestOnlineUsersInterval = setInterval(requestOnlineUsers, 10_000)
        intervals.add(requestOnlineUsersInterval)
    }

    static drawOnlineUsers(users) {
        Array.from(document.getElementsByClassName("chat")).forEach(el => { el.classList.remove('online') })

        if (users != null) {
            const user = Utils.getUser()

            const onlineUsersEl = document.getElementById("online-users");
            if (onlineUsersEl) {
                onlineUsersEl.innerText = ""
                users.sort((a, b) => a.firstName < b.firstName ? 1 : -1)
                users.forEach((u) => {
                    if (u.id == user.id) {
                        return
                    }

                    var chat = document.getElementById(`chat-${u.id}`)
                    if (!chat) {
                        chat = newChatElement({ user: u })
                        onlineUsersEl.prepend(chat)
                    }
                    chat.classList.add('online')
                })
            }
        }
    }

    static drawChats(chats) {
        if (chats != null) {
            const chatsEl = document.getElementById("users-chats");
            chatsEl.innerHTML = ""
            chats.forEach((chat) => {
                const el = newChatElement(chat)
                chatsEl.append(el)
            })
        }
    }

    static async drawNewMessage(message) {
        const user = Utils.getUser()

        if ((message.senderID == recipientID || message.senderID == user.id) && !(document.getElementById(`message-${message.id}`))) {
            const chatMessages = document.getElementById("chat-messages");
            const el = newMessageElement(message)
            chatMessages.appendChild(el)
            chatMessages.scrollTop = chatMessages.scrollHeight - chatMessages.clientHeight;
        }

        var chatId = message.senderID == user.id ? message.recipientID : message.senderID
        const chat = document.getElementById(`chat-${chatId}`)
        if (chat) {
            if (message.recipientID == user.id ) {
                changeChatUnreadCount(document.getElementById(`chat-${chatId}-unread-messages-count`), 1)
            }
            document.getElementById(`chat-${chatId}-lastMessage`).innerText = message.message
            document.getElementById(`chat-${chatId}-lastMessageDate`).innerText = `${new Date(message.date).toLocaleString()}`
            const chatsEl = document.getElementById("users-chats");
            chatsEl.prepend(chat)
        } else {
            requestChats()
        }
    }

    static async prependMessages(messages) {
        const chatMessages = document.getElementById("chat-messages");
        const scrollToEnd = (chatMessages.childNodes.length == 0)

        if (messages == null) {
            chatMessages.removeEventListener("scroll", loadMessages)
            return
        }

        messages.forEach((message) => {
            const el = newMessageElement(message)
            chatMessages.prepend(el)

            if (scrollToEnd) {
                chatMessages.scrollTop = chatMessages.scrollHeight
            }
        })
    }

    static async changeMessageStatusToRead(messageID) {
        setTimeout(() => {
            let el = document.getElementById(`message-${messageID}-status`)
            if (el) {
                el.innerText = '✓✓'
            }
        }, 1000)
    }
}