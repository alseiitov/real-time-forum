import fetcher from "../services/Fetcher.js";
import AbstractView from "./AbstractView.js";

const genders = { 1: 'Male', 2: 'Female' }

const getUserByID = async (id) => {
    const path = `/api/users/${id}`

    return await fetcher.get(path);
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Profile");
        this.userID = params.userID;
    }

    async getHtml() {
        return `
            <div id="user-profile">
                <div class="profile-info" id="avatar"></div>
                <div class="profile-info" id="username"></div>
                <div class="profile-info" id="first-name"></div>
                <div class="profile-info" id="last-name"></div>
                <div class="profile-info" id="age"></div>
                <div class="profile-info" id="gender"></div>
                <div class="profile-info" id="role"></div>
                <div class="profile-info" id="registered"></div>
            </div>
        `;
    }

    async init() {
        const user = await getUserByID(this.userID)

        document.querySelector('.profile-info#avatar').innerHTML = `<img src="http://${API_HOST_NAME}/images/${user.avatar}">`
        document.querySelector('.profile-info#username').innerText = `Username: ${user.username}`
        document.querySelector('.profile-info#first-name').innerText = `First name: ${user.firstName}`
        document.querySelector('.profile-info#last-name').innerText = `Last name: ${user.lastName}`
        document.querySelector('.profile-info#age').innerText = `Age: ${user.age}`
        document.querySelector('.profile-info#gender').innerText = `Gender: ${genders[user.gender]}`
        document.querySelector('.profile-info#role').innerText = `Role: ${user.role}`
        document.querySelector('.profile-info#registered').innerText = `Registered: ${new Date(Date.parse(user.registered)).toLocaleString()}`
    }
}