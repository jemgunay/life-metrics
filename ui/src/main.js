import "bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";
import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import DayLogPage from "./components/pages/DayLogPage.vue";
import QueryPage from "./components/pages/QueryPage.vue";

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
];

const router = createRouter({
    history: createWebHistory(),
    routes: routes
});

createApp(App).use(router).mount("#app");
