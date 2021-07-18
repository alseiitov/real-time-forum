import intervals from "../services/Intervals.js";
import Ws from "../services/Ws.js";
import AbstractView from "./AbstractView.js";
import Utils from "../services/Utils.js"

var requestOnlineUsersInterval
var recipientID

var loadMessages
var sendTypingInEvent

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

    if (!message.read) {
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

    const typingIndicator = document.createElement("div")
    typingIndicator.classList.add('typing-indicator')
    typingIndicator.id = `chat-${chat.user.id}-typing-indicator`
    typingIndicator.innerHTML = `<span></span><span></span><span></span>`
    messageEl.append(typingIndicator)

    const lastMessage = document.createElement("p")
    lastMessage.id = `chat-${chat.user.id}-lastMessage`

    const lastMessageDate = document.createElement("p")
    lastMessageDate.id = `chat-${chat.user.id}-lastMessageDate`

    if (chat.lastMessage.id) {
        lastMessage.innerText = `${chat.lastMessage.message}`
        lastMessageDate.innerText = `${new Date(chat.lastMessage.date).toLocaleString()}`
    }

    messageEl.append(lastMessage)
    messageEl.append(lastMessageDate)

    const unreadMessagesCount = document.createElement("div")
    unreadMessagesCount.classList.add(`unread-messages-count`)
    unreadMessagesCount.id = `chat-${chat.user.id}-unread-messages-count`
    changeChatUnreadCount(unreadMessagesCount, chat.unreadMessagesCount)

    el.append(messageEl)
    el.append(unreadMessagesCount)

    el.addEventListener("click", () => {
        Array.from(document.getElementsByClassName("chat")).forEach(el => { el.classList.remove('active') })
        Array.from(document.getElementsByClassName("unread-messages-count")).forEach(el => { el.classList.remove('active') })

        el.classList.add('active')
        unreadMessagesCount.classList.add('active')

        document.getElementById("message-form").style.display = "block"
        document.getElementById("message-input").value = ""
        document.getElementById("chat-messages").innerHTML = ""
        document.getElementById("chat-messages").addEventListener("scroll", loadMessages)
        recipientID = chat.user.id
        updateQueryParams()
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

const updateQueryParams = () => {
    const urlParams = new URLSearchParams(window.location.search)
    urlParams.set('user', recipientID)
    history.replaceState(null, null, "?" + urlParams.toString())
}

const newTypingInListener = () => {
    return Utils.throttle(() => {
        Ws.send(JSON.stringify({ type: "typingInRequest", body: { recipientID: recipientID } }))
    }, 2000)
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Chats");
    }

    async getHtml() {
        return `
            <div id="chats-container">
                <div id="chats"></div>
                <div>
                    <div id="chat-messages">
                        <p id="chat-messages-placeholder">Select chat to start messaging</p>
                    </div>
                    <form id="message-form">
                        <input type="text" id="message-input" size="64" placeholder="Send message" autocomplete="off" autofocus/>
                    </form>
                </div>
            </div>
        `;
    }



    async init() {
        const urlParams = new URLSearchParams(window.location.search)
        recipientID = parseInt(urlParams.get('user')) || 0

        const chatMessages = document.getElementById("chat-messages");
        const messageForm = document.getElementById("message-form");
        messageForm.style.display = 'none'

        
        const messageInput = document.getElementById("message-input");
        sendTypingInEvent = newTypingInListener()
        messageInput.addEventListener("input", sendTypingInEvent)

        loadMessages = Utils.debounce(function () {
            if (chatMessages.scrollTop < chatMessages.scrollHeight * 0.1) {
                let offsetMsg = document.querySelector('.message')
                if (!offsetMsg) {
                    chatMessages.removeEventListener("scroll", loadMessages)
                    return
                }
                let offsetMsgID = parseInt(offsetMsg.id.split('-')[1]);
                Ws.send(JSON.stringify({ type: "messagesRequest", body: { userID: recipientID, lastMessageID: offsetMsgID } }))
            }

            if (chatMessages.scrollTop == 0) {
                chatMessages.scrollTop = 10
            }
        }, 100)

        chatMessages.addEventListener("scroll", loadMessages)
        chatMessages.addEventListener("scroll", () => {

        })

        messageForm.onsubmit = function () {
            if (!messageInput.value) {
                return false;
            }

            let message = { recipientID: recipientID, message: messageInput.value }
            Ws.send(JSON.stringify({ type: "message", body: message }));
            messageInput.value = "";

            messageInput.removeEventListener("input", sendTypingInEvent)
            sendTypingInEvent = newTypingInListener()
            messageInput.addEventListener("input", sendTypingInEvent)

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
            users.forEach((u) => {
                var chat = document.getElementById(`chat-${u.id}`)
                if (chat) {
                    chat.classList.add('online')
                } else {
                    const user = Utils.getUser()
                    if (u.id != user.id) {
                        requestChats()
                        requestOnlineUsers()
                    }
                }
            })
        }
    }

    static drawChats(chats) {
        if (chats != null) {
            const chatsEl = document.getElementById("chats");
            if (chatsEl) {
                chatsEl.innerHTML = ""
                chats.forEach((chat) => {
                    const el = newChatElement(chat)
                    chatsEl.append(el)
                    if (recipientID == chat.user.id) {
                        el.click()
                        el.scrollIntoView({ behavior: "smooth" })
                    }
                })
            }
        }
    }

    static async drawNewMessage(message) {
        const user = Utils.getUser()

        if ((message.senderID == recipientID || (message.senderID == user.id && message.recipientID == recipientID)) && !(document.getElementById(`message-${message.id}`))) {
            const chatMessages = document.getElementById("chat-messages");
            if (chatMessages) {
                const el = newMessageElement(message)
                chatMessages.appendChild(el)
                chatMessages.scrollTop = chatMessages.scrollHeight - chatMessages.clientHeight;
            }
        }

        var chatId = message.senderID == user.id ? message.recipientID : message.senderID
        const chat = document.getElementById(`chat-${chatId}`)
        if (chat) {
            if (message.recipientID == user.id) {
                changeChatUnreadCount(document.getElementById(`chat-${chatId}-unread-messages-count`), 1)
            }

            document.getElementById(`chat-${chatId}-lastMessage`).innerText = message.message
            document.getElementById(`chat-${chatId}-lastMessage`).style.display = ""

            document.getElementById(`chat-${chatId}-lastMessageDate`).innerText = `${new Date(message.date).toLocaleString()}`

            document.getElementById(`chat-${chatId}-typing-indicator`).classList.remove("typing")
            document.getElementById(`chat-${chatId}-typing-indicator`).style.display = "none"

            const chatsEl = document.getElementById("chats");
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

    static async changeMessageStatusToRead(message) {
        const user = Utils.getUser()
        if (message.recipientID == user.id) {
            changeChatUnreadCount(document.getElementById(`chat-${message.senderID}-unread-messages-count`), -1)
        }

        setTimeout(() => {
            let el = document.getElementById(`message-${message.id}-status`)
            if (el) {
                el.innerText = '✓✓'
            }
        }, 1000)
    }

    static async drawTypingIn(event) {
        const indicator = document.getElementById(`chat-${event.senderID}-typing-indicator`)
        const lastMessage = document.getElementById(`chat-${event.senderID}-lastMessage`)

        if (indicator) {
            indicator.style.display = "table"
            lastMessage.style.display = "none"

            indicator.classList.add("typing")
            setTimeout(() => {
                indicator.classList.remove("typing")
            }, 2000)

            setTimeout(() => {
                if (!indicator.classList.contains("typing")) {
                    indicator.style.display = "none"
                    lastMessage.style.display = ""
                }
            }, 3000)
        }
    }
}