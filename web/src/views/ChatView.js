import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Chat");
        this.userID = params.userID
    }

    async getHtml() {
        return `
            <div id="log"></div>
            <form id="form">
                <input type="submit" value="Send" />
                <input type="text" id="msg" size="64" autofocus />
            </form>
        `;
    }

    async init() {
        var conn;
        var msg = document.getElementById("msg");
        var log = document.getElementById("log");

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

        const loadMessages = debounce(function () {
            console.log(log.scrollHeight)
            console.log(log.scrollTop)

            if (log.scrollTop < log.scrollHeight * 0.1) {
                let offsetMsg = document.querySelector('.message')
                let offsetMsgID = parseInt(offsetMsg.id.split('-')[1]);
                conn.send(JSON.stringify({ type: "messagesRequest", body: { userID: recipientID, lastMessageID: offsetMsgID } }))
            }

            if (log.scrollTop == 0) {
                log.scrollTop = 1
            }
        }, 50)

        log.addEventListener("scroll", loadMessages)

        var recipientID = parseInt(this.userID)
    
        if (recipientID == 1) {
            var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiI5Mzg4MDQ0NTY5Iiwicm9sZSI6IjIiLCJzdWIiOiIyIn0.dLVytdQkOTxiNAvia4GfEjk6IJGvdHlygrjBCKFm9KU"
            var senderID = 2

        }
        if (recipientID == 2) {
            var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiI5Mzg4MDQ0NTQ3Iiwicm9sZSI6IjIiLCJzdWIiOiIxIn0.2bRxGzoFv2MNEWixcl-ZFIvh3NSlTxkLcV49q5wnOMY"
            var senderID = 1
        }

        function appendLog(item) {
            var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
            log.appendChild(item);
            if (doScroll) {
                log.scrollTop = log.scrollHeight - log.clientHeight;
            }
        }

        document.getElementById("form").onsubmit = function () {
            if (!conn) {
                return false;
            }
            if (!msg.value) {
                return false;
            }

            conn.send(JSON.stringify({ type: "message", body: { recipientID: recipientID, message: msg.value } }));
            msg.value = "";
            return false;
        };

        if (window["WebSocket"]) {
            conn = new WebSocket("ws://127.0.0.1:8081/ws")
            conn.onopen = function (evt) {
                conn.send(JSON.stringify({ type: "token", body: token }))
                conn.send(JSON.stringify({ type: "messagesRequest", body: { userID: recipientID, lastMessageID: 0 } }))
            }

            conn.onerror = function (evt) {
                alert(evt.data)
            }
            conn.onclose = function (evt) {
                var item = document.createElement("div");
                item.innerHTML = "<b>Connection closed.</b>";
                appendLog(item);
            };
            conn.onmessage = function (evt) {
                let obj = JSON.parse(evt.data)
                console.log(obj)
                switch (obj.type) {
                    case "message":
                        let msg = obj.body
                        var item = document.createElement("div");
                        item.id = `message-${msg.id}`

                        if (!msg.read && msg.recipientID == senderID) {
                            conn.send(JSON.stringify({ type: "readMessageRequest", body: { messageID: msg.id } }))
                        }
                        if (!msg.read && msg.senderID == senderID) {
                            item.style.color = "gray"
                        }

                        item.classList.add("message")
                        item.innerText = `user ${msg.senderID} sends to user ${msg.recipientID}\n${new Date(Date.parse(msg.date)).toLocaleString()}\n${msg.message}\n`
                        appendLog(item);

                        break
                    case "messagesResponse":
                        let scrollToEnd = log.childNodes.length == 0

                        obj.body.forEach((msg) => {
                            var item = document.createElement("div");
                            item.id = `message-${msg.id}`

                            if (!msg.read && msg.recipientID == senderID) {
                                conn.send(JSON.stringify({ type: "readMessageRequest", body: { messageID: msg.id } }))
                            }
                            if (!msg.read && msg.senderID == senderID) {
                                item.style.color = "gray"
                            }

                            item.classList.add("message")
                            item.innerText = `user ${msg.senderID} sends to user ${msg.recipientID}\n${new Date(Date.parse(msg.date)).toLocaleString()}\n${msg.message}\n`

                            log.prepend(item)

                            if (scrollToEnd) {
                                log.scrollTop = log.scrollHeight
                            }
                        })

                        break
                    case "readMessageResponse":

                        let el = document.getElementById(`message-${obj.body}`)
                        el.style.color = "black"

                    case "notification":
                        // console.log(obj)
                        break
                    case "error":
                        alert(obj.body)
                        break
                    case "pingMessage":
                        conn.send(JSON.stringify({ type: "pongMessage" }))
                        break

                }
            };
        } else {
            var item = document.createElement("div");
            item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
            appendLog(item);
        }
    }
}