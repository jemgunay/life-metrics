<template>
    <nav class="navbar navbar-expand-md navbar-light bg-light flex-nowrap mb-3">
        <!-- collapse burger icon -->
        <button class="navbar-toggler mr-3" type="button" data-toggle="collapse" data-target="#navbar-main">
            <span class="navbar-toggler-icon"></span>
        </button>
        <!-- brand -->
        <router-link class="navbar-brand w-100" active-class="active" :to="{ name: 'home'}">Life Log</router-link>
        <!-- navs -->
        <div class="navbar-collapse collapse w-100 justify-content-center" id="navbar-main">
            <ul class="navbar-nav mx-auto">
                <li class="nav-item">
                    <router-link class="nav-link" active-class="active" :to="{ name: 'home'}">Day Log</router-link>
                </li>
                <li class="nav-item">
                    <router-link class="nav-link" active-class="active" :to="{ name: 'sources'}">Sources</router-link>
                </li>
            </ul>
        </div>
        <!-- offset to centre navs -->
        <div class="w-100"></div>
    </nav>
</template>

<script>
export default {
    name: "NavBar",
    data() {
        return {};
    }
};
</script>

<style scoped>
.navbar {
    border-bottom: 1px solid rgba(0,0,0,.125);
}
</style>