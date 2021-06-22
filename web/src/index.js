import NavBar from "./views/NavBarView.js";
import Home from "./views/HomeView.js";
import SignUp from "./views/SignUpView.js";
import SignIn from "./views/SignInView.js";
import Chats from "./views/ChatsView.js";
import Chat from "./views/ChatView.js";

import Error404 from "./views/error404View.js";
import Error500 from "./views/error500View.js";
import Error503 from "./views/error503View.js";

import Ws from "./services/Ws.js"
import Utils from "./services/Utils.js"

const pathToRegex = path => new RegExp("^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "(.+)") + "$");

const getParams = match => {
    const values = match.result.slice(1);
    const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(result => result[1]);

    return Object.fromEntries(keys.map((key, i) => {
        return [key, values[i]];
    }));
};

const navigateTo = url => {
    history.pushState(null, null, url);
    router();
};

const router = async () => {
    const routes = [
        { path: "/", view: Home },
        { path: "/sign-up", view: SignUp },
        { path: "/sign-in", view: SignIn },
        { path: "/chats", view: Chats },
        { path: "/chat/:userID", view: Chat },
        { path: "/404", view: Error404 },
        { path: "/500", view: Error500 },
        { path: "/503", view: Error503 },
    ];

    // Test each route for potential match
    const potentialMatches = routes.map(route => {
        return {
            route: route,
            result: location.pathname.match(pathToRegex(route.path))
        };
    });

    let match = potentialMatches.find(potentialMatch => potentialMatch.result !== null);

    if (!match) {
        match = {
            route: routes[0],
            result: [location.pathname]
        };
    }

    const user = Utils.getUser()

    const navBarView = new NavBar(null, user);
    const pageView = new match.route.view(getParams(match), user);

    document.querySelector("#navbar").innerHTML = await navBarView.getHtml();
    document.querySelector("#app").innerHTML = await pageView.getHtml();
    
    pageView.init()
};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", async () => {
    document.body.addEventListener("click", function (e) {
        if (e.target.matches("[data-link]")) {
            e.preventDefault();
            navigateTo(e.target.href);
        }
    });

    if (localStorage.getItem("accessToken")) {
        await Ws.connect()
    }

    router()
});


export default { navigateTo };