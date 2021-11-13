import "bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";
import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import DayLogPage from "./components/pages/DayLogPage.vue";
import SourcesPage from "./components/pages/SourcesPage.vue";

const routes = [
    {
        path: "/",
        name: "home",
        component: DayLogPage
    },
    {
        path: "/sources",
        name: "sources",
        component: SourcesPage
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes: routes
});

import { library } from "@fortawesome/fontawesome-svg-core";
import { faCheckCircle, faTimesCircle } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

library.add(faCheckCircle, faTimesCircle);

createApp(App).use(router).component("font-awesome-icon", FontAwesomeIcon).mount("#app");

