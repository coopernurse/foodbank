const m = window.m;

const HouseholdService = {
    createHousehold: async (household) => {
        try {
            const response = await m.request({
                method: "POST",
                url: "/household",
                body: household,
                withCredentials: true,
            });
            return response;
        } catch (err) {
            console.error("Failed to create household:", err);
            throw err;
        }
    },
};

export default HouseholdService;
