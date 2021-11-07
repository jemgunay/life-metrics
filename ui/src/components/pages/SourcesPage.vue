<template>
    <div class="row">
        <div class="col-md-12">
            <div class="alert alert-dismissible fade show" v-bind:class="'alert-' + alertIndicator"
                 role="alert" v-if="alertIndicator">
                <span>{{ alertMessage }}</span>

                <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
        </div>

        <div class="col-md-12 mb-3">
            <div class="card">
                <h5 class="card-header d-flex justify-content-between align-items-center">
                    Sources
                    <button type="button" class="btn btn-sm btn-success" v-on:click="performCollectRequest">Collect
                    </button>
                </h5>
                <div class="card-body">
                    <p class="mb-0">View source state and configuration below.</p>
                </div>
            </div>
        </div>

        <div class="col-md-6">
            <div class="card">
                <h5 class="card-header">Monzo</h5>
                <div class="card-body">
                    <a :href="apiHost + '/api/auth/monzo'" target="_blank">Click to authenticate
                        Monzo.</a>
                    <p class="mb-0">Ensure to approve the email and the app notification.</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from "axios";

export default {
    name: "SourcesPage",
    data() {
        return {
            apiHost: process.env.VUE_APP_API_HOST,
            alertIndicator: "",
            alertMessage: "",

        };
    },
    methods: {
        setBanner(state, msg) {
            this.alertIndicator = state;
            this.alertMessage = msg;
        },

        performCollectRequest() {
            this.setBanner();

            axios({
                method: "POST",
                url: process.env.VUE_APP_API_HOST + "/api/data/collect",
                validateStatus: function (status) {
                    return status === 202;
                }
            })
                .then(() => {
                    this.setBanner("success", "Successfully submitted a source collection request!");
                })
                .catch((error) => {
                    this.setBanner("danger", "Source collection request failed! " + error);
                    console.error(error);
                });
        }
    }
};
</script>

<style>
</style>