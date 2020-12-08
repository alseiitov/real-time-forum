var container = document.querySelector(".container");

fetch('http://localhost:8081/api')
    .then((response) => {
        return response.json();
    })
    .then((data) => {
        container.textContent = JSON.stringify(data)
    });