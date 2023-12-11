import Login from "../components/Login.vue";
import Home from "../components/Home.vue";
import { createRouter, createWebHashHistory } from "vue-router";

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home
    },
    {
        path: '/login',
        name: 'Login',
        component: Login
    }
];

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router;