import AuthState from "./services/AuthState.js";

const m = window.m;

const Shell = {
  view: (vnode) => {
    return m("div", { class: "min-h-screen bg-white" }, [
      // Navbar
      m("nav", {
        class: "bg-gray-600 p-4 text-white flex items-center justify-between",
      }, [
        m("img", {
          src: "/static/img/mend-logo.png",  // Replace with your actual logo image URL
          alt: "Food Bank Logo",
          class: "h-10"  // Adjust size of the logo
        }),
        m("div", [
          m("a", {
            href: "#",
            class: "text-white hover:text-gray-200 mx-4"
          }, "Dashboard"),
          m("a", {
            href: "#",
            class: "text-white hover:text-gray-200 mx-4"
          }, "Food Banks"),
          m("a", {
            href: "#",
            class: "text-white hover:text-gray-200 mx-4"
          }, "Visits"),
          m("a", {
            href: "#",
            class: "text-white hover:text-gray-200 mx-4"
          }, "Items"),
          AuthState.isAuthenticated() && m("a", {
            href: "#",
            class: "text-white hover:text-gray-200 mx-4",
            onclick: () => AuthState.logout()
          }, "Logout"),
        ])
      ]),

      // Main content area
      m("main", {
        class: "p-6 bg-gray-100"
      }, vnode.children)  // Dynamic content based on what's passed into the shell
    ]);
  }
};

export default Shell;
