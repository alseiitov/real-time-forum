import fetcher from "../services/Fetcher.js";
import AbstractView from "./AbstractView.js";

const genders = { 1: 'Male', 2: 'Female' }

const getUserByID = async (id) => {
    const path = `/api/users/${id}`
    return await fetcher.get(path);
}

const getUsersPosts = async (userID) => {
    const path = `/api/users/${userID}/posts`
    return await fetcher.get(path);
}

const getUsersRatedPosts = async (userID) => {
    const path = `/api/users/${userID}/rated-posts`
    return await fetcher.get(path);
}

const newPostElement = (post) => {
    const el = document.createElement("div")
    el.classList.add("post")

    const linkToPost = document.createElement("a")
    linkToPost.setAttribute("href", `/post/${post.id}`)
    linkToPost.setAttribute("data-link", "")
    linkToPost.innerText = `${post.title}`

    const postDate = document.createElement("p")
    postDate.innerText = new Date(post.date).toLocaleString()

    const linkToAuthor = document.createElement("a")
    linkToAuthor.setAttribute("href", `/users/${post.author.id}`)
    linkToAuthor.setAttribute("data-link", "")
    linkToAuthor.innerText = `${post.author.firstName} ${post.author.lastName}`

    el.append(linkToPost)
    el.append(postDate)
    el.append(linkToAuthor)

    return el
}

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Profile");
        this.userID = params.userID;
    }

    async getHtml() {
        return `
        <h2>Users profile</h2>
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
            <h2>Users posts</h2>
            <div id="users-posts"></div>
            <h2>Users liked posts</h2>
            <div id="users-liked-posts"></div>
            <h2>Users disliked posts</h2>
            <div id="users-disliked-posts"></div>
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

        const usersPosts = await getUsersPosts(this.userID)
        const usersRatedPosts = await getUsersRatedPosts(this.userID)

        const usersPostsEl = document.getElementById('users-posts')
        if (usersPosts != null) {
            usersPosts.forEach((post) => {
                const postEl = newPostElement(post)
                usersPostsEl.append(postEl)
            })
        } else {
            usersPostsEl.innerText = 'No posts'
        }
     
        
        const usersLikedPosts = usersRatedPosts.filter((post) => post.userRate == 1)
        const usersLikedPostsEl = document.getElementById('users-liked-posts')
        if(usersLikedPosts != null) {
            usersLikedPosts.forEach((post) => {
                const postEl = newPostElement(post)
                usersLikedPostsEl.append(postEl)
            })
        } else {
            usersLikedPostsEl.innerText = 'No posts'
        }

        const usersDisLikedPosts = usersRatedPosts.filter((post) => post.userRate == 2)
        const usersDislikedPostsEl = document.getElementById('users-disliked-posts')
        if(usersDisLikedPosts != null) {
            usersDisLikedPosts.forEach((post) => {
                const postEl = newPostElement(post)
                usersDislikedPostsEl.append(postEl)
            })
        } else {
            usersDislikedPostsEl.innerText = 'No posts'
        }
    }
}