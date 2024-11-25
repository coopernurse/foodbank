const m = window.m;
import AuthState from "../services/AuthState.js";

export const Login = {
    oninit: () => {
        Login.email = "";
        Login.password = "";
        Login.error = null;
    },
    email: "",
    password: "",
    error: null,
    onSubmit() {
        Login.error = null; // Clear any existing error messages
        const email = Login.email;
        const password = Login.password;
        AuthState.login(email, password)
            .then(() => {
                m.route.set("/"); // Redirect to home on success
            })
            .catch((err) => {
                Login.error = "Invalid email or password"; // Set error message
                Login.password = ""; // Clear the password field
                console.error("Login failed:", err);
            });
    },
    view() {
        return m(
            "div",
            {
                class: "flex items-center justify-center min-h-screen bg-gray-100",
            },
            m(
                "div",
                {
                    class: "absolute top-[30%] w-full max-w-md px-6 py-8 bg-white shadow-md rounded-lg",
                },
                [
                    m("h1", { class: "text-2xl font-bold mb-6 text-center" }, "Login"),
                    Login.error &&
                        m("div", {
                            class: "mb-4 text-sm text-red-600",
                        }, Login.error),
                    m("form", { class: "space-y-4", onsubmit: (e) => e.preventDefault() }, [
                        m("div", [
                            m("label", { class: "block text-sm font-medium text-gray-700 mb-1" }, "Email"),
                            m("input", {
                                class:
                                    "w-full px-3 py-2 border rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500 border-gray-300",
                                type: "email",
                                placeholder: "you@example.com",
                                oninput: (e) => (Login.email = e.target.value),
                                value: Login.email,
                            }),
                        ]),
                        m("div", [
                            m("label", { class: "block text-sm font-medium text-gray-700 mb-1" }, "Password"),
                            m("input", {
                                class:
                                    "w-full px-3 py-2 border rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500 border-gray-300",
                                type: "password",
                                placeholder: "••••••••",
                                oninput: (e) => (Login.password = e.target.value),
                                value: Login.password,
                            }),
                        ]),
                        m(
                            "button",
                            {
                                class:
                                    "w-full py-2 px-4 text-white bg-blue-600 rounded-lg hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2",
                                type: "submit",
                                onclick: Login.onSubmit,
                            },
                            "Login"
                        ),
                    ]),
                ]
            )
        );
    },
};

export default Login;
