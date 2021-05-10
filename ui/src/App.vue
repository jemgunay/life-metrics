<template>
    <div>
        <NavBar/>
        <div class="container">
            <div class="row">
                <router-view></router-view>
            </div>
        </div>
    </div>
</template>

<script>
import NavBar from "./components/NavBar.vue";
//import Monzo from "./components/Monzo.vue";

export default {
    name: "App",
    components: {
        NavBar,
        //Monzo
    }
};
</script>

<style>
</style>
