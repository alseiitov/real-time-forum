var container = document.querySelector(".container");

var APIAdress = `http://${APIAdress}/api/user`

fetch(APIAdress)
    .then((response) => {
        return response.json();
    })
    .then((data) => {
        container.textContent = JSON.stringify(data)
    });
