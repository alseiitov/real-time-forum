var container = document.querySelector(".container");

var APIAdress = `http://${config.API.host}:${config.API.port}/api`

fetch(APIAdress)
    .then((response) => {
        return response.json();
    })
    .then((data) => {
        container.textContent = JSON.stringify(data)
    });
