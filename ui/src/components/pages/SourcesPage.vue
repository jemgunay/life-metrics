<template>
    <div class="row">
        <div class="col-md-12 mb-3">
            <div class="card">
                <h5 class="card-header">Sources</h5>
                <div class="card-body">
                    <p class="mb-0">View source state and configuration below.</p>
                </div>
            </div>
        </div>

        <div class="col-md-6">
            <div class="card">
                <h5 class="card-header">Monzo</h5>
                <div class="card-body">
                    <a href="http://localhost:8080/api/auth/monzo">Click to authenticate Monzo.</a>
                    <p class="mb-0">Ensure you approve the email and in the app.</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "SourcesPage",
    data() {
        return {};
    }
};
</script>

<style>
</style>