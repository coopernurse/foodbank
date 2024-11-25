const m = window.m;

export const ResetPassword = {
    oninit: () => {
        ResetPassword.email = "";
        ResetPassword.error = null;
    },
    email: "",
    error: null,
    onSubmit() {
        ResetPassword.error = null; // Clear any existing error messages
        const email = ResetPassword.email;
        m.request({
            method: "POST",
            url: "/send-password-reset-email",
            body: { email },
            withCredentials: true, // In case cookies are used for auth
        }).then(() => {
            ResetPassword.error = "Password reset email sent. Please check your inbox.";
        }).catch((err) => {
            ResetPassword.error = "Failed to send password reset email. Please try again.";
            console.error("Password reset failed:", err);
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
                    m("h1", { class: "text-2xl font-bold mb-6 text-center" }, "Reset Password"),
                    ResetPassword.error &&
                        m("div", {
                            class: "mb-4 p-3 bg-red-400 rounded-lg text-gray-800 font-bold",
                        }, ResetPassword.error),
                    m("form", { class: "space-y-4", onsubmit: (e) => e.preventDefault() }, [
                        m("div", [
                            m("label", { class: "block text-sm font-medium text-gray-700 mb-1" }, "Email"),
                            m("input", {
                                class:
                                    "w-full px-3 py-2 border rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500 border-gray-300",
                                type: "email",
                                placeholder: "you@example.com",
                                oninput: (e) => (ResetPassword.email = e.target.value),
                                value: ResetPassword.email,
                            }),
                        ]),
                        m(
                            "button",
                            {
                                class:
                                    "w-full py-2 px-4 text-white bg-blue-600 rounded-lg hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2",
                                type: "submit",
                                onclick: ResetPassword.onSubmit,
                            },
                            "Send Reset Email"
                        ),
                    ]),
                ]
            )
        );
    },
};

export default ResetPassword;
