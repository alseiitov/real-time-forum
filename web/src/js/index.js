var container = document.querySelector(".container");

var APIAdress = `http://${APIAdress}/api/users/1`

fetch(APIAdress)
    .then((response) => {
        return response.json();
    })
    .then((data) => {
        container.textContent = JSON.stringify(data)
    });
