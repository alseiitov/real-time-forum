import AbstractView from "./AbstractView.js";
import fetcher from "../services/Fetcher.js"
import router from "../index.js"

var currCategoryID = 1
var currPageNum = 1
var postsEnded = false

const getCategories = async () => {
    const path = "/api/categories"

    const response = await fetcher.get(path)
    if (response.ok) {
        const data = await response.json()
        return data
    }
}

const drawCategories = async (categories) => {
    const categoriesEl = document.getElementById("categories")
    categories.forEach(category => {
        const el = document.createElement("button")
        el.innerText = category.name
        el.id = `category-${category.id}`

        el.addEventListener("click", async () => {
            const titleEl = document.getElementById("category-title")
            currCategoryID = category.id
            currPageNum = 1
            postsEnded = false
            titleEl.innerText = category.name
            document.getElementById("page-number").innerText = currPageNum

            drawPostsByCategoryID(category.id, currPageNum)
        })

        categoriesEl.append(el)
    })
}

const drawPostsByCategoryID = async (categoryID, page) => {
    const postsEl = document.getElementById("posts")
    postsEl.innerHTML = ""

    const path = `/api/categories/${categoryID}/${page}`

    const response = await fetcher.get(path)
    switch (response.status) {
        case 200:
            const data = await response.json()
            if (data.posts) {
                console.log(data.posts)
                data.posts.forEach((post) => {
                    const postEl = newPostElement(post)
                    postsEl.append(postEl)
                })
            } else {
                postsEl.innerText = "No posts"
                postsEnded = true
            }
            break
        case 404:
            router.navigateTo("/404")
            break
    }
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
        this.setTitle("Home");
    }

    async getHtml() {
        return `
            <h1>Hello</h1>

            <div id="categories"></div>
            <div id="category-title"></div>
           
            <div id="posts"></div>
            <button id="prev-button">prev</button>
            <p id="page-number"></p>
            <button id="next-button">next</button>
        `;
    }

    async init() {
        const categories = await getCategories()
        drawCategories(categories)

        document.getElementById(`category-${currCategoryID}`).click()

        const nextButtonEl = document.getElementById(`next-button`)
        const prevButtonEl = document.getElementById(`prev-button`)
        const pageNumber = document.getElementById(`page-number`)

        nextButtonEl.addEventListener("click", () => {
            if (postsEnded) {
                return
            }
            currPageNum++
            pageNumber.innerText = currPageNum

            drawPostsByCategoryID(currCategoryID, currPageNum)
        })

        prevButtonEl.addEventListener("click", () => {
            if (currPageNum == 1) {
                return
            }
            postsEnded = false
            currPageNum--
            pageNumber.innerText = currPageNum

            drawPostsByCategoryID(currCategoryID, currPageNum)
        })
    }
}