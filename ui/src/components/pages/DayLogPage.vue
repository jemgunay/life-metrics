<template>
    <div class="col-md-12">
        <div class="card">
            <h5 class="card-header">Day Log</h5>
            <div class="card-body">
                <form class="row">
                    <div class="col-md-12">
                        <div class="alert alert-dismissible fade show" v-bind:class="'alert-' + alertIndicator" role="alert" v-if="alertIndicator">
                            <span v-if="alertIndicator === 'success'">Day log submitted!</span>
                            <span v-else>Day log submission failed!</span>
                            <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                                <span aria-hidden="true">&times;</span>
                            </button>
                        </div>
                    </div>

                    <div class="form-group col-md-6">
                        <label for="log-date-input">Log Date</label>
                        <input type="date" class="form-control" id="log-date-input" v-model="logDate">
                    </div>

                    <div class="form-group col-md-6">
                        <label for="general-mood-input">General Mood</label>
                        <input type="range" class="form-control" id="general-mood-input" min="0" max="10"
                               v-model="logMetrics['general_mood']">
                    </div>

                    <div class="form-group col-md-6">
                        <label for="work-mood-input">Work Mood</label>
                        <input type="range" class="form-control" id="work-mood-input" min="0" max="10"
                               v-model="logMetrics['work_mood']">
                    </div>

                    <div class="form-group col-md-6">
                        <label for="diet-input">Diet Quality</label>
                        <input type="range" class="form-control" id="diet-input" min="0" max="10"
                               v-model="logMetrics['diet_quality']">
                    </div>

                    <div class="form-group col-md-6">
                        <label for="water-input">Water Intake</label>
                        <input type="range" class="form-control" id="water-input" min="0" max="10"
                               v-model="logMetrics['water_intake']">
                    </div>

                    <div class="form-group col-md-6">
                        <label for="sleep-input">Sleep Quality</label>
                        <input type="range" class="form-control" id="sleep-input" min="0" max="10"
                               v-model="logMetrics['sleep_quality']">
                    </div>

                    <div class="form-group col-6 col-md-3">
                        <div class="custom-control custom-checkbox">
                            <input type="checkbox" class="custom-control-input" id="exercise-check"
                                   v-model="logMetrics['exercise']">
                            <label class="custom-control-label" for="exercise-check">Exercise</label>
                        </div>
                    </div>

                    <div class="form-group col-6 col-md-3">
                        <div class="custom-control custom-checkbox">
                            <input type="checkbox" class="custom-control-input" id="meditation-check"
                                   v-model="logMetrics['meditation']">
                            <label class="custom-control-label" for="meditation-check">Meditation</label>
                        </div>
                    </div>

                    <div class="form-group col-md-12">
                        <label for="notes-input">Notes</label>
                        <textarea class="form-control" id="notes-input" v-model="logNotes"></textarea>
                    </div>

                    <div class="form-group col-md-12 text-right mb-0">
                        <button type="submit" class="btn btn-primary" v-on:click="submitDayLog">Submit</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<script>
import axios from "axios";

export default {
    name: "DayLogPage",
    data() {
        return {
            alertIndicator: "",
            logDate: (new Date()).toISOString().slice(0,10),
            logNotes: "",
            logMetrics: {},
        };
    },
    methods: {
        submitDayLog: function(event) {
            event.preventDefault();
            this.alertIndicator = "";

            console.log(this.logMetrics);
            let reqBody = {
                "date": this.logDate + "T00:00:00Z",
                "notes": this.logNotes,
                "metrics": {...this.logMetrics},
            };
            console.log(reqBody);

            axios({
                method: "post",
                url: "http://localhost:8080/api",
                data: JSON.stringify(reqBody)
            })
                .then((resp) => {
                    this.alertIndicator = "success";
                    console.log(resp);
                })
                .catch((error) => {
                    this.alertIndicator = "danger";
                    console.error(error);
                });
        }
    }
};
</script>

<style>
</style>