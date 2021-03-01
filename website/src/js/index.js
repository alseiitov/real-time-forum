const http = {
    StatusOK: 200,
    StatusCreated: 201,
    StatusBadRequest: 400,
    StatusInternalServerError: 500
}

const parseImageInBase64 = async () => {
    let image = document.querySelector('input[name="image"]').files[0]
    if (!image) {
        return ""
    }
    return await fileToBase64(image)
}

const fileToBase64 = (file) => {
    return new Promise(resolve => {
        let fileReader = new FileReader();

        fileReader.onload = (fileLoadedEvent) => {
            resolve(fileLoadedEvent.target.result)
        }

        fileReader.readAsDataURL(file)
    })
}

const checkImage = (base64) => /image\/(jpeg|png|gif)/.test(base64)


const submitButton = document.getElementById('submit-button')

submitButton.addEventListener('click', async () => {
    let base64 = await parseImageInBase64()
    console.log(base64)
    console.log(checkImage(base64))
})

