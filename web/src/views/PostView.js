import AbstractView from "./AbstractView.js";
import fetcher from "../services/Fetcher.js"
import router from "../index.js"
import Utils from "../services/Utils.js";

const likeTypes = {
    like: 1,
    dislike: 2
}

var currCommentPageNum = 1
var commentsEnded = false

const getPost = async (postID) => {
    const path = `/api/posts/${postID}`
    return await fetcher.get(path)
}

const getComments = async (postID, page) => {
    const path = `/api/posts/${postID}/comments/${page}`
    return await fetcher.get(path)
}

const addComment = async (postID, data, image) => {
    const path = `/api/posts/${postID}/comments`
    const body = { data: data, image: image }

    fetcher.post(path, body)
}


const likePost = async (postID, likeType) => {
    const path = `/api/posts/${postID}/likes`
    const body = { likeType: likeType }

    fetcher.post(path, body)

    const likeButton = document.getElementById(`like-post-button`)
    const dislikeButton = document.getElementById(`dislike-post-button`)
    const rating = document.getElementById("post-rating")

    const alreadyLiked = likeButton.classList.contains('rated')
    const alreadyDisliked = dislikeButton.classList.contains('rated')

    likeButton.classList.remove('rated')
    dislikeButton.classList.remove('rated')

    if (likeType == likeTypes.like) {
        if (alreadyLiked) {
            rating.innerText--
        } else {
            likeButton.classList.add('rated')
            rating.innerText++
        }
        if (alreadyDisliked) {
            rating.innerText++
        }
    }


    if (likeType == likeTypes.dislike) {
        if (alreadyDisliked) {
            rating.innerText++
        } else {
            dislikeButton.classList.add('rated')
            rating.innerText--
        }
        if (alreadyLiked) {
            rating.innerText--
        }
    }
}

const likeComment = async (commentID, likeType) => {
    const path = `/api/comments/${commentID}/likes`
    const body = { likeType: likeType }

    await fetcher.post(path, body)

    const likeButton = document.getElementById(`like-comment-${commentID}`)
    const dislikeButton = document.getElementById(`dislike-comment-${commentID}`)

    const alreadyLiked = likeButton.classList.contains('rated')
    const alreadyDisliked = dislikeButton.classList.contains('rated')
    const rating = document.getElementById(`comment-${commentID}-rating`)


    likeButton.classList.remove('rated')
    dislikeButton.classList.remove('rated')


    if (likeType == likeTypes.like) {
        if (alreadyLiked) {
            rating.innerText--
        } else {
            likeButton.classList.add('rated')
            rating.innerText++
        }
        if (alreadyDisliked) {
            rating.innerText++
        }
    }


    if (likeType == likeTypes.dislike) {
        if (alreadyDisliked) {
            rating.innerText++
        } else {
            dislikeButton.classList.add('rated')
            rating.innerText--
        }
        if (alreadyLiked) {
            rating.innerText--
        }
    }
}

const drawPost = async (post) => {
    document.getElementById("post-title").innerText = post.title;

    if (post.image) {
        document.getElementById("post-image").innerHTML = `<img src="http://${API_HOST_NAME}/images/${post.image}">`;
    }

    document.getElementById("post-data").innerText = post.data;

    const postAuthor = document.createElement("a")
    postAuthor.innerText = `${post.author.firstName} ${post.author.lastName}`
    postAuthor.setAttribute("href", `/user/${post.author.id}`)
    postAuthor.setAttribute("data-link", "")
    document.getElementById("post-author").append(postAuthor)

    document.getElementById("post-creation-date").innerText = new Date(post.date).toLocaleString();

    const likeButton = document.getElementById("like-post-button")
    if (likeButton) {
        likeButton.addEventListener("click", () => { likePost(post.id, likeTypes.like) })
        if (post.userRate == likeTypes.like) {
            likeButton.classList.add('rated')
        }
    }

    const dislikeButton = document.getElementById("dislike-post-button")
    if (dislikeButton) {
        dislikeButton.addEventListener("click", () => { likePost(post.id, likeTypes.dislike) })
        if (post.userRate == likeTypes.dislike) {
            dislikeButton.classList.add('rated')
        }
    }

    document.getElementById("post-rating").innerText = post.rating;

    const categoriesEl = document.getElementById("post-categories")

    for (let i = 1; i < post.categories.length; i++) {
        const category = post.categories[i];

        const categoryEl = document.createElement("a")
        categoryEl.innerText = `${category.name}`
        categoryEl.setAttribute("href", `/?category=${category.id}&page=1}`)
        categoryEl.setAttribute("data-link", "")

        categoriesEl.append(categoryEl)
    }
}

const drawPostCommentsPage = async (postID, page) => {
    const data = await getComments(postID, page)
    drawPostComments(data)
}

const drawPostComments = async (comments) => {
    const commentsEl = document.getElementById("post-comments")
    if (!comments) {
        commentsEnded = true
        commentsEl.innerText = "No comments"
        return
    }

    commentsEl.innerText = ""

    comments.forEach(comment => { drawComment(comment, false) })
}

const drawComment = (comment, isNewComment) => {
    const user = Utils.getUser()
    const commentsEl = document.getElementById("post-comments")

    const commentEl = document.createElement("div")
    commentEl.classList.add("post-comment")

    const commentAuthor = document.createElement("a")
    if (comment.author.id == user.id) {
        commentAuthor.innerText = `You`
    } else {
        commentAuthor.innerText = `${comment.author.firstName} ${comment.author.lastName}`
    }

    commentAuthor.setAttribute("href", `/user/${comment.author.id}`)
    commentAuthor.setAttribute("data-link", "")
    commentEl.append(commentAuthor)

    if (comment.image) {
        const commentImage = document.createElement("img")
        commentImage.src = `http://${API_HOST_NAME}/images/${comment.image}`
        commentEl.append(commentImage)
    }

    const commentText = document.createElement("p")
    commentText.innerText = `${comment.data}\n${new Date(comment.date).toLocaleString()}`
    commentEl.append(commentText)

    const likeButton = document.createElement("button")
    likeButton.classList.add("rate-button")
    likeButton.id = `like-comment-${comment.id}`
    likeButton.innerText = "▴"
    likeButton.addEventListener("click", () => { likeComment(comment.id, likeTypes.like) })
    if (comment.userRate == likeTypes.like) {
        likeButton.classList.add('rated')
    }
    commentEl.append(likeButton)

    const dislikeButton = document.createElement("button")
    dislikeButton.classList.add("rate-button")
    dislikeButton.id = `dislike-comment-${comment.id}`
    dislikeButton.innerText = "▾"
    dislikeButton.addEventListener("click", () => { likeComment(comment.id, likeTypes.dislike) })
    if (comment.userRate == likeTypes.dislike) {
        dislikeButton.classList.add('rated')
    }
    commentEl.append(dislikeButton)

    const rating = document.createElement("p")
    rating.id = `comment-${comment.id}-rating`
    rating.innerText = comment.rating
    commentEl.append(rating)


    if (isNewComment) {
        commentsEl.prepend(commentEl)
    } else {
        commentsEl.append(commentEl)
    }
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
                <button class="rate-button" id="like-post-button">▴</button>
                <button class="rate-button" id="dislike-post-button">▾</button>
                `
                :
                `<p>Sign-in to rate a post</p>`
            )
            +
            `
                <p id="post-rating"></p>
            </div>
            <div id="post-categories"></div>
            `
            +
            (authorized ?
                `
                <div id="post-comments"></div>
                <button id="prev-button">prev</button>
                <p id="page-number">1</p>
                <button id="next-button">next</button>
                <form id="comment-form" onsubmit="return false;">
                    <textarea id="comment-input" cols="30" rows="5" maxlength="128" placeholder="Leave a comment" required></textarea><br>
                    <input type="file" id="comment-image-input" accept="image/jpeg, image/png, image/gif">
                    <div id="comment-image-preview"></div>
                    <button>Send</button>
                    <div id="input-error"></div>
                </form>

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
        if (post) {
            drawPost(post, this.user)

            const authorized = Boolean(this.user.id)
            if (authorized) {
                drawPostCommentsPage(this.postID, currCommentPageNum)

                const nextButtonEl = document.getElementById(`next-button`)
                const prevButtonEl = document.getElementById(`prev-button`)
                const pageNumber = document.getElementById(`page-number`)

                nextButtonEl.addEventListener("click", () => {
                    if (commentsEnded) {
                        return
                    }
                    currCommentPageNum++
                    pageNumber.innerText = currCommentPageNum
                    drawPostCommentsPage(this.postID, currCommentPageNum)
                })

                prevButtonEl.addEventListener("click", () => {
                    if (currCommentPageNum == 1) {
                        return
                    }
                    commentsEnded = false
                    currCommentPageNum--
                    pageNumber.innerText = currCommentPageNum
                    drawPostCommentsPage(this.postID, currCommentPageNum)
                })
            }

            const commentText = document.getElementById("comment-input")
            const imageInput = document.getElementById("comment-image-input")
            const imagePreview = document.getElementById("comment-image-preview")
            const inputError = document.getElementById("input-error")

            const imageMaxSize = 20 * 1024 * 1024
            const allowedImageTypes = ["image/jpeg", "image/png", "image/gif"]

            var imageBase64 = ""

            imageInput.addEventListener("change", async () => {
                inputError.innerText = ""
                imagePreview.innerHTML = ""

                const image = imageInput.files[0]
                if (image.size > imageMaxSize) {
                    inputError.innerText = "Too big image! Max image size is 20 Mb"
                    imageInput.value = ""
                    imageBase64 = ""
                    return
                }

                if (!allowedImageTypes.includes(image.type)) {
                    inputError.innerText = `Only ${allowedImageTypes.join(", ")} types are allowed`
                    imageInput.value = ""
                    imageBase64 = ""
                    return
                }

                imageBase64 = await Utils.fileToBase64(image)
                imagePreview.innerHTML = `<img src="${imageBase64}">`
            })


            document.getElementById("comment-form").addEventListener("submit", async () => {
                const comment = await addComment(this.postID, commentText.value, imageBase64)
                drawComment(comment, true)

                imageInput.value = ""
                commentText.value = ""
            })
        }
    }
}

