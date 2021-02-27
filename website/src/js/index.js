var container = document.querySelector(".container");

const http = {
    StatusOK: 200,
    StatusCreated: 201,
    StatusBadRequest: 400,
    StatusInternalServerError: 500
}

console.log(APIAdress)

let savePhoto = async (inp) => {
    let image = inp.files[0];
    if (!image) return ""

    let formData = new FormData();
    formData.append("image", image);

    const options = {
        method: "POST",
        body: formData,
    }

    try {
        const resp = await fetch('http://localhost:8081/api/images', options)
        const obj = await resp.json()

        if (resp.status != http.StatusCreated) {
            alert(resp.status + obj.error)
        } else {
            return obj.name
        }
    } catch (e) {
        alert(`Can't upload image, error: ${e.message}`);
    }
}


var submitButton = document.getElementById('submit-button')

submitButton.addEventListener('click', async () => {
    let imageInput = document.querySelector('input[name="image"]')
    let name = await savePhoto(imageInput)
    console.log(name)
})

