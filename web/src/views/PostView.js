import AbstractView from "./AbstractView.js";
import fetcher from "../services/Fetcher.js"
import router from "../index.js"

const likeTypes = {
    like: 1,
    dislike: 2
}

const getPost = async (postID) => {
    const path = `/api/posts/${postID}`

    const response = await fetcher.get(path)
    if (response.ok) {
        return await response.json()
    }
}

const getComments = async (postID, page) => {
    const path = `/api/posts/${postID}/comments/${page}`

    const response = await fetcher.get(path)
    switch (response.status) {
        case 200:
            const data = await response.json()
            return data
        case 400:
            router.navigateTo("/400")
            break
    }
}

const likeComment = async (commentID, likeType) => {
    const path = `/api/comments/${commentID}/likes`
    const body = { likeType: likeType }

    const response = await fetcher.post(path, body)
    if (response.status == 400) {
        router.navigateTo("/400")
        return
    }

    const likeButton = document.getElementById(`like-comment-${commentID}`)
    const dislikeButton = document.getElementById(`dislike-comment-${commentID}`)

    const alreadyLiked = likeButton.classList.contains('rated')
    const alreadyDisliked = dislikeButton.classList.contains('rated')

    likeButton.classList.remove('rated')
    dislikeButton.classList.remove('rated')

    if (likeType == likeTypes.like && !alreadyLiked) {
        likeButton.classList.add('rated')
    }

    if (likeType == likeTypes.dislike && !alreadyDisliked) {
        dislikeButton.classList.add('rated')
    }
}

const drawPost = async (post) => {
    document.getElementById("post-title").innerText = post.title;

    if (post.image) {
        document.getElementById("post-image").innerHTML = `<img src="http://${API_HOST_NAME}/images/${post.image}">`;
    }

    document.getElementById("post-data").innerText = post.data;
    //TODO: parse user name
    document.getElementById("post-author").innerText = post.userID;
    document.getElementById("post-creation-date").innerText = new Date(post.date).toLocaleString();

    const categoriesEl = document.getElementById("post-categories")
    post.categories.forEach(category => {
        const el = document.createElement("a")
        el.innerText = category.name
        el.href = `/category/${category.id}/1`
        categoriesEl.append(el)
    })
}


const drawPostComments = async (comments, userID) => {
    const commentsEl = document.getElementById("post-comments")
    comments.forEach(comment => {
        const commentEl = document.createElement("div")
        commentEl.classList.add("post-comment")

        const commentAuthor = document.createElement("a")
        //TODO: parse user name
        commentAuthor.innerText = `user ${comment.userID}`
        commentAuthor.setAttribute("href", `/user/${comment.userID}`)
        commentAuthor.setAttribute("data-link", "")

        const commentText = document.createElement("p")
        commentText.innerText = `${comment.data}\n${new Date(comment.date).toLocaleString()}`

        const likeButton = document.createElement("button")
        likeButton.classList.add("rate-button")
        likeButton.id = `like-comment-${comment.id}`
        likeButton.innerText = "like"
        likeButton.addEventListener("click", () => { likeComment(comment.id, likeTypes.like) })

        const dislikeButton = document.createElement("button")
        dislikeButton.classList.add("rate-button")
        dislikeButton.id = `dislike-comment-${comment.id}`
        dislikeButton.innerText = "dislike"
        dislikeButton.addEventListener("click", () => { likeComment(comment.id, likeTypes.dislike) })

        commentEl.append(commentAuthor)
        commentEl.append(commentText)
        commentEl.append(likeButton)
        commentEl.append(dislikeButton)
        commentsEl.append(commentEl)
    })
}

export default class extends AbstractView {
    constructor(params, user) {
        super(params);
        this.setTitle("Post");
        this.user = user
        this.postID = params.postID
    }

    async getHtml() {
        const authorized = Boolean(this.user.id)

        return `
            <div class="post">
                <div id="post-title"></div>
                <div id="post-image"></div>
                <div id="post-data"></div>
                <div id="post-info">
                    <div id="post-author"></div>
                    <div id="post-creation-date"></div>
                </div>
                <div id="likes">
            `
            +
            (authorized ?
                `
                <button class="rate-button" id="like-post-button">like</button>
                <button class="rate-button" id="dislike-post-button">dislike</button>
                `
                :
                `<p>Sign-in to rate a post</p>`
            )
            +
            `
                <p id="likes-count">0</p>
            </div>
            <div id="post-categories"></div>
            `
            +
            (authorized ?
                `
                <div id="post-comments"></div>
                `
                :
                `
                <p>Sign-in to read or leave a comment</p>
                `
            )
            +
            `
            </div>
            `
    }

    async init() {
        const post = await getPost(this.postID)
        drawPost(post, this.user)

        const authorized = Boolean(this.user.id)
        if (authorized) {
            const data = await getComments(this.postID, 1)
            if (data) {
                drawPostComments(data, this.userID)
            }
        }
    }
}