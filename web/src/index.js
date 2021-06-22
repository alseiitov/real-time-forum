import NavBar from "./views/NavBarView.js";
import Home from "./views/HomeView.js";
import SignUp from "./views/SignUpView.js";
import SignIn from "./views/SignInView.js";
import Chats from "./views/ChatsView.js";
import Chat from "./views/ChatView.js";

import Error401 from "./views/error401View.js";
import Error404 from "./views/error404View.js";
import Error500 from "./views/error500View.js";
import Error503 from "./views/error503View.js";

import Ws from "./services/Ws.js"
import Utils from "./services/Utils.js"

const roles = {
    guest: 1,
    user: 2,
    moderator: 3,
    admins: 4,
}

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
        { path: "/", view: Home, minRole: roles.guest },
        { path: "/sign-up", view: SignUp, minRole: roles.guest },
        { path: "/sign-in", view: SignIn, minRole: roles.guest },
        { path: "/chats", view: Chats, minRole: roles.user },
        { path: "/chat/:userID", view: Chat, minRole: roles.user },
        { path: "/401", view: Error401, minRole: roles.guest },
        { path: "/404", view: Error404, minRole: roles.guest },
        { path: "/500", view: Error500, minRole: roles.guest },
        { path: "/503", view: Error503, minRole: roles.guest },
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
    if (!user.role) {
        user.role = roles.guest
        localStorage.setItem('role', user.role)
    }

    if (user.role < match.route.minRole) {
        navigateTo("/401")
        return
    }

    const navBarView = new NavBar(null, user);
    const pageView = new match.route.view(getParams(match), user);

    document.querySelector("#navbar").innerHTML = await navBarView.getHtml();
    document.querySelector("#app").innerHTML = await pageView.getHtml();

    pageView.init()
    navBarView.init()
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