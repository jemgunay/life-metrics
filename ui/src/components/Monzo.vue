<template>
    <div class="col-md-12">
        <h1>Howdy</h1>

        <form>
            <div class="form-group">
                <input type="hidden" class="form-control" name="client_id"/>
            </div>
            <div class="form-group">
                <input type="hidden" class="form-control" name="redirect_uri"/>
            </div>
            <div class="form-group">
                <input type="hidden" class="form-control" name="response_type"/>
            </div>
            <button type="button" class="btn btn-primary">Authorise</button>
        </form>
    </div>
</template>

<script>
export default {
    name: "Monzo",
    data() {
        return {};
    }
};
</script>

<style scoped>
</style>
