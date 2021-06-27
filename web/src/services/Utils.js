const parseJwt = (token) => {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
};

const getUser = () => {
    return {
        id: localStorage.getItem('sub'),
        role: localStorage.getItem('role'),
        accessToken: localStorage.getItem('accessToken'),
        refreshToken: localStorage.getItem('refreshToken')
    }
}

const logOut = () => {
    localStorage.removeItem('sub')
    localStorage.removeItem('role')
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
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

const base64isImage = (base64string) => {
    return /image\/(jpeg|png|gif)/.test(base64string)
}


export default { parseJwt, getUser, logOut, fileToBase64, base64isImage }