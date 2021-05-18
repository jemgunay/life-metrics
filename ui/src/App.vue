<template>
    <div>
        <NavBar/>
        <div class="container">
            <div class="row">
                <DayLogPage/>
            </div>
        </div>
    </div>
</template>

<script>
import NavBar from "./components/NavBar.vue";
import DayLogPage from "./components/pages/DayLogPage.vue";

export default {
    name: "App",
    components: {
        NavBar,
        DayLogPage
    }
};
</script>

<style>
</style>
