import AbstractView from "./AbstractView.js";
import fetcher from "../services/Fetcher.js"
import router from "../index.js"
import Utils from "../services/Utils.js"

const getCategories = async () => {
    const path = "/api/categories"
    return await fetcher.get(path)
}


const createPost = async (post) => {
    const path = `/api/posts`
    return await fetcher.post(path, post)
}


export default class extends AbstractView {
    constructor(params, user) {
        super(params);
        this.setTitle("New post");
        this.user = user
    }

    async getHtml() {
        return `
            <form id="post-form" onsubmit="return false;">
                <p>Title:</p>
                <input type="text" id="title-input" minlength="2" maxlength="64" required>

                <p>Image:</p>
                <label for="post-image-input" class="custom-file-input">
                    Choose image
                </label>
                <input type="file" id="post-image-input" accept="image/jpeg, image/png, image/gif">
                <div id="post-image-preview"></div>

                <p>Categories:</p>
                <select class="category-selector" id="category-1">
                    <option value="0" selected disabled>Category</option>
                </select>
                <select class="category-selector" id="category-2">
                    <option value="0" selected disabled>Category</option>
                </select>
                <select class="category-selector" id="category-3">
                    <option value="0" selected disabled>Category</option>
                </select>

                <textarea id="data-input" cols="30" rows="5" minlength="2" maxlength="512" required></textarea>
                <div class="error" id="error-message"></div>
                <button>Send</button>
            </form>
        `;
    }

    async init() {
        const categories = await getCategories()

        Array.from(document.getElementsByClassName("category-selector")).forEach(el => {
            for (let i = 1; i < categories.length; i++) {
                const optionEl = document.createElement("option")
                optionEl.value = categories[i].id;
                optionEl.innerText = categories[i].name
                el.append(optionEl)
            }
        })

        const postTitle = document.getElementById("title-input")
        const imageInput = document.getElementById("post-image-input")
        const imagePreview = document.getElementById("post-image-preview")
        const category1 = document.getElementById("category-1")
        const category2 = document.getElementById("category-2")
        const category3 = document.getElementById("category-3")
        const postData = document.getElementById("data-input")
        const errorMessage = document.getElementById("error-message")

        const imageMaxSize = 20 * 1024 * 1024
        const allowedImageTypes = ["image/jpeg", "image/png", "image/gif"]

        var imageBase64 = ""

        imageInput.addEventListener("change", async () => {
            errorMessage.innerText = ""
            imagePreview.innerHTML = ""

            const image = imageInput.files[0]
            if (image.size > imageMaxSize) {
                errorMessage.innerText = "Too big image! Max image size is 20 Mb"
                imageInput.value = ""
                imageBase64 = ""
                return
            }

            if (!allowedImageTypes.includes(image.type)) {
                errorMessage.innerText = `Only ${allowedImageTypes.join(", ")} types are allowed`
                imageInput.value = ""
                imageBase64 = ""
                return
            }

            imageBase64 = await Utils.fileToBase64(image)
            imagePreview.innerHTML = `<img src="${imageBase64}">`
        })

        document.getElementById("post-form").addEventListener("submit", async () => {
            errorMessage.innerText = ""

            const categories = [parseInt(category1.value), parseInt(category2.value), parseInt(category3.value)].filter(n => n != 0)
            if (!categories.length) {
                errorMessage.innerText = "Please select at least one category"
                return
            }

            const post = await createPost(
                {
                    title: postTitle.value,
                    image: imageBase64,
                    categories: categories,
                    data: postData.value
                }
            )
            
            if (post) {
                router.navigateTo(`/post/${post.postID}`)
            }
        })
    }
}



