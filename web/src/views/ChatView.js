import AbstractView from "./AbstractView.js";
import Ws from "../services/Ws.js"
import Utils from "../services/Utils.js"


var loadMessages

const newMessageElement = (message) => {
    const el = document.createElement("div");
    el.id = `message-${message.id}`

    const userID = parseInt(localStorage.getItem("sub"))

    if (!message.read && message.recipientID == userID) {
        Ws.send(JSON.stringify({ type: "readMessageRequest", body: { messageID: message.id } }))
    }

    if (!message.read && message.senderID == userID) {
        el.classList.add("unread")
    }

    el.classList.add("message")
    if (message.senderID == userID) {
        el.classList.add("sended-message")
    } else {
        el.classList.add("received-message")
    }

    const messageText = document.createElement('p')
    messageText.classList.add('message-text')
    messageText.innerText = message.message

    const messageDate = document.createElement('p')
    messageDate.classList.add('message-date')
    messageDate.innerText = new Date(Date.parse(message.date)).toLocaleString()

    el.append(messageText)
    el.append(messageDate)

    return el
}

const debounce = (func, wait, immediate) => {
    var timeout;
    return function () {
        var context = this, args = arguments;
        var later = function () {
            timeout = null;
            if (!immediate) func.apply(context, args);
        };
        var callNow = immediate && !timeout;
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
        if (callNow) func.apply(context, args);
    };
};




export default class extends AbstractView {
    constructor(params, user) {
        super(params);
        this.setTitle("Chat");
        this.user = user
        this.recipientID = params.userID
    }

    async getHtml() {
        return `
            <div id="chat-messages"></div>
            <form id="message-form">
                <input type="text" id="message-input" size="64" autofocus />
            </form>
        `;
    }

    async init() {
        const chatMessages = document.getElementById("chat-messages");
        const messageForm = document.getElementById("message-form");
        const messageInput = document.getElementById("message-input");
        const recipientID = parseInt(this.recipientID)

        loadMessages = debounce(function () {
            if (chatMessages.scrollTop < chatMessages.scrollHeight * 0.1) {
                let offsetMsg = document.querySelector('.message')
                let offsetMsgID = parseInt(offsetMsg.id.split('-')[1]);
                Ws.send(JSON.stringify({ type: "messagesRequest", body: { userID: recipientID, lastMessageID: offsetMsgID } }))
            }

            if (chatMessages.scrollTop == 0) {
                chatMessages.scrollTop = 1
            }
        }, 100)


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

        Ws.send(JSON.stringify({ type: "messagesRequest", body: { userID: recipientID, lastMessageID: 0 } }))
    }

    static async appendNewMessage(message) {
        const user = Utils.getUser()
        const currChatUserID = location.pathname.split("/").pop()
        if ((message.senderID == currChatUserID || message.senderID == user.id) && !(document.getElementById(`message-${message.id}`))) {
            const chatMessages = document.getElementById("chat-messages");
            const el = newMessageElement(message)
            chatMessages.appendChild(el)
            chatMessages.scrollTop = chatMessages.scrollHeight - chatMessages.clientHeight;
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
}



