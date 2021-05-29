import "bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";
import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import DayLogPage from "./components/pages/DayLogPage.vue";
import QueryPage from "./components/pages/QueryPage.vue";
import ConfigPage from "./components/pages/ConfigPage.vue";

const routes = [
    {
        path: "/daylog",
        name: "dayLog",
        component: DayLogPage,
    },
    {
        path: "/",
        name: "home",
        redirect: {name: "dayLog"}
    },
    {
        path: "/query",
        name: "query",
        component: QueryPage,
    },
    {
        path: "/config",
        name: "config",
        component: ConfigPage,
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes: routes
});


import { library } from "@fortawesome/fontawesome-svg-core";
import { faInfoCircle } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

library.add(faInfoCircle);

createApp(App).use(router).component("font-awesome-icon", FontAwesomeIcon).mount("#app");

