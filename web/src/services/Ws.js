import Chat from "../views/ChatView.js";
import Chats from "../views/ChatsView.js"

var connection

const getConnection = () => {
    return new Promise((resolve, reject) => {
        if (connection != undefined && connection.readyState) {
            resolve(connection)
            return
        }

        if (window["WebSocket"]) {
            let token = localStorage.getItem("accessToken")

            if (token == undefined) {
                alert("error opening websocket connection, no access token in localStorage")
                return
            }

            let conn = new WebSocket(`ws://${API_HOST_NAME}/ws`)

            conn.onerror = function (evt) {
                alert(evt.data)
                reject(evt)
            }

            conn.onclose = function (evt) {
                alert("connection closed")
            };

            conn.onmessage = function (evt) {
                let obj = JSON.parse(evt.data)
                console.log(obj)
                switch (obj.type) {
                    case "message":
                        Chat.appendNewMessage(obj.body)
                        break
                    case "messagesResponse":
                        Chat.prependMessages(obj.body)
                        break
                    case "readMessageResponse":
                        let el = document.getElementById(`message-${obj.body}`)
                        el.style.color = ""

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
    }

}

export default Ws

