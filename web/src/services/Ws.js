import Chat from "../views/ChatView.js";
import Chats from "../views/ChatsView.js"
import Utils from "./Utils.js";

var connection

const getConnection = () => {
    if (connection && connection.readyState < 2) {
        return Promise.resolve(connection)
    }

    return new Promise((resolve, reject) => {
        if (window["WebSocket"]) {
            let token = localStorage.getItem("accessToken")

            if (token == undefined) {
                alert("error opening websocket connection, no access token in localStorage")
                return
            }

            const conn = new WebSocket(`ws://${API_HOST_NAME}/ws`)

            conn.onerror = function (evt) {
                Utils.showError(503)
                return
            }

            conn.onmessage = function (evt) {
                let obj = JSON.parse(evt.data)

                switch (obj.type) {
                    case "message":
                        switch (true) {
                            case location.pathname.startsWith("/chat/"):
                                Chat.appendNewMessage(obj.body)
                                break
                            case location.pathname == "/chats":
                                conn.send(JSON.stringify({ type: "chatsRequest" }))
                                break
                        }
                        break
                    case "messagesResponse":
                        Chat.prependMessages(obj.body)
                        break
                    case "chatsResponse":
                        Chats.drawChats(obj.body)
                        break
                    case "readMessageResponse":
                        let el = document.getElementById(`message-${obj.body}`)
                        if (el) {
                            el.classList.remove('unread')
                        }

                    case "notification":
                        // console.log(obj)
                        break
                    case "onlineUsersResponse":
                        Chats.drawOnlineUsers(obj.body)
                        break
                    case "error":
                        alert(obj.body)
                        break
                    case "pingMessage":
                        conn.send(JSON.stringify({ type: "pongMessage" }))
                        break
                }
            };

            conn.onopen = function () {
                conn.send(JSON.stringify({ type: "token", body: token }))
                resolve(conn)
            }
        } else {
            alert("Your browser does not support WebSockets")
        }
    })
}


const Ws = {
    connect: async () => {
        connection = await getConnection()
    },

    send: async (e) => {
        connection = await getConnection()
        connection.send(e)
    },

    disconnect: async () => {
        connection.close()
    }
}

export default Ws

