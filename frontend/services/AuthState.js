const m = window.m;

const AuthState = {
    token: null,
    isAuthenticated: () => AuthState.token !== null,

    login: async (email, password) => {
        try {
            const response = await m.request({
                method: "POST",
                url: "/login",
                body: { email, password },
                withCredentials: true, // In case cookies are used for auth
            });

            AuthState.token = response.sessionToken;
            if (!AuthState.token) {
                throw new Error("Authentication failed");
            }
        } catch (err) {
            console.error("Failed to authenticate:", err);
            throw err; // Propagate the error
        }
    },

    logout: () => {
        AuthState.token = null;
        m.route.set("/login"); // Redirect to login page after logout
    },
};

export default AuthState;
