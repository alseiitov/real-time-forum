import AbstractView from "./AbstractView.js";
import Ws from "../services/Ws.js"


var loadMessages

const newMessageElement = (message) => {
    const el = document.createElement("div");
    el.id = `message-${message.id}`

    const userID = parseInt(localStorage.getItem("sub"))

    if (!message.read && message.recipientID == userID) {
        Ws.send(JSON.stringify({ type: "readMessageRequest", body: { messageID: message.id } }))
    }

    if (!message.read && message.senderID == userID) {
        el.style.color = "gray"
    }

    el.classList.add("message")
    el.innerText = `user ${message.senderID} sends to user ${message.recipientID}\n${new Date(Date.parse(message.date)).toLocaleString()}\n${message.message}\n`

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
    constructor(params) {
        super(params);
        this.setTitle("Chat");
        this.userID = params.userID
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
        const recipientID = parseInt(this.userID)

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
        const chatMessages = document.getElementById("chat-messages");
        const doScroll = chatMessages.scrollTop > chatMessages.scrollHeight - chatMessages.clientHeight - 1;
        const el = newMessageElement(message)
        chatMessages.appendChild(el);
        if (doScroll) {
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



