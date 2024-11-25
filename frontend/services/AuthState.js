const m = window.m;

const AuthState = {
    token: null,
    isAuthenticated: () => AuthState.token !== null,

    login: async (email, password) => {
        try {
            const response = await m.request({
                method: "POST",
                url: "/api/login",
                body: { email, password },
                withCredentials: true, // In case cookies are used for auth
            });

            AuthState.token = response.token;
            return true; // Login succeeded
        } catch (err) {
            console.error("Failed to authenticate:", err);
            return false; // Login failed
        }
    },
};

export default AuthState;
