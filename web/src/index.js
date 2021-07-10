import NavBar from "./views/NavBarView.js";
import Home from "./views/HomeView.js";
import SignUp from "./views/SignUpView.js";
import SignIn from "./views/SignInView.js";
import Post from "./views/PostView.js";
import NewPost from "./views/NewPostView.js";
import Chats from "./views/ChatsView.js";
import Profile from "./views/ProfileView.js";

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
        { path: "/new-post", view: NewPost, minRole: roles.user },
        { path: "/post/:postID", view: Post, minRole: roles.guest },
        { path: "/chats", view: Chats, minRole: roles.user },
        { path: "/user/:userID", view: Profile, minRole: roles.user },
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
        Utils.showError(404)
        return
    }

    const user = Utils.getUser()
    if (!user.role) {
        user.role = roles.guest
        localStorage.setItem('role', user.role)
    }

    if (user.role < match.route.minRole) {
        Utils.showError(401)
        return
    }

    const navBarView = new NavBar(null, user);
    const pageView = new match.route.view(getParams(match), user);

    document.querySelector("#navbar").innerHTML = await navBarView.getHtml();
    navBarView.init()

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