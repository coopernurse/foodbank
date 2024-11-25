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
            return AuthState.token != ""; // Login succeeded
        } catch (err) {
            console.error("Failed to authenticate:", err);
            return false; // Login failed
        }
    },

    logout: () => {
        AuthState.token = null;
        m.route.set("/login"); // Redirect to login page after logout
    },
};

export default AuthState;
