import m from 'mithril';
import Login from "./auth/Login";
import Signup from "./auth/Signup";
import Shell from "./Shell";
import AuthState from "./services/AuthState";

import ResetPassword from "./auth/ResetPassword";

m.route(document.body, "/login", {
    "/login": {
        render: () => {
            if (AuthState.isAuthenticated()) {
                m.route.set("/"); // Redirect to the main app if logged in
            }
            return m(Shell, m(Login));
        },
    },
    "/signup": {
        render: () => {
            return m(Shell, m(Signup));
        },
    },
    "/reset-password": {
        render: () => {
            return m(Shell, m(ResetPassword));
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
