const m = window.m;
import Login from "./auth/Login";
import Shell from "./Shell";
import AuthState from "./services/AuthState";

m.route(document.body, "/login", {
    "/login": {
        render: () => {
            if (AuthState.isAuthenticated()) {
                m.route.set("/"); // Redirect to the main app if logged in
            }
            return m(Shell, m(Login));
        },
    },
    "/": {
        render: () => {
            if (!AuthState.isAuthenticated()) {
                m.route.set("/login"); // Redirect to login if not authenticated
            }
            return m(Shell, "Welcome to the Food Bank Management App!"); // Placeholder for the main page
        },
    },
});
