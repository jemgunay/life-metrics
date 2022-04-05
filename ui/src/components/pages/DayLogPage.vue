<template>
    <div class="row">
        <div class="col-md-12">
            <div class="card">
                <h5 class="card-header">Day Log</h5>
                <div class="card-body">
                    <form class="row">
                        <div class="col-md-12">
                            <div class="alert alert-dismissible fade show" v-bind:class="'alert-' + alertIndicator"
                                 role="alert" v-if="alertIndicator">
                                <span>{{ alertMessage }}</span>

                                <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                                    <span aria-hidden="true">&times;</span>
                                </button>
                            </div>
                        </div>


                        <div class="form-group col-md-6">
                            <label for="log-date-input">Log Date</label>
                            <div class="input-group mb-3">
                                <input type="date" class="form-control" id="log-date-input" v-model="logDate" v-on:change="getDayLog">
                                <div class="input-group-append">
                                    <button class="btn btn-outline-secondary" type="button" v-on:click="resetDate">Today</button>
                                </div>
                            </div>
                        </div>

                        <div class="form-group col-md-6">
                            <label for="general-mood-input">General Mood</label>
                            <input type="range" class="form-control" id="general-mood-input" min="0" max="10"
                                   v-model="logMetrics['general_mood']">
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
                            <label for="caffeine-input">Caffeine Intake</label>
                            <input type="range" class="form-control" id="caffeine-input" min="0" max="10"
                                   v-model="logMetrics['caffeine_intake']">
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
    </div>
</template>

<script>
import axios from "axios";

export default {
    name: "DayLogPage",
    data() {
        return {
            alertIndicator: "",
            alertMessage: "",
            logDate: "",
            logMetrics: {
                "general_mood": 5,
                "diet_quality": 5,
                "water_intake": 5,
                "caffeine_intake": 1,
                "exercise": false,
                "meditation": false
            },
            logMetricsDefaults: {},
            logNotes: ""
        };
    },
    mounted() {
        // store a copy of the metrics defaults for form resetting
        this.logMetricsDefaults = this.logMetrics;
        // determine if log has been submitted today already
        this.resetDate();
    },
    methods: {
        setBanner(state, msg) {
            this.alertIndicator = state;
            this.alertMessage = msg;
        },

        resetDate() {
            let newLogDate = (new Date()).toISOString().slice(0, 10);
            // only get day log data if the date has changed
            if (newLogDate === this.logDate) {
                return;
            }
            this.logDate = newLogDate;
            this.getDayLog();
        },

        getDayLog() {
            let date = this.logDate + "T00:00:00Z";

            this.performDayLogRequest("/daylog/data/daylog?date=" + date, "GET", "", (data) => {
                if (data["submitted"] === true) {
                    this.logMetrics = data["metrics"];
                    // invert caffeine (lower the better)
                    this.logMetrics["caffeine_intake"] = 10 - this.logMetrics["caffeine_intake"];
                    this.logNotes = data["notes"];

                    this.setBanner("success", "Day log completed for the selected day.");
                    return;
                }

                // reset form defaults
                this.logMetrics = this.logMetricsDefaults;
                this.logNotes = "";
                this.setBanner("info", "Day log not completed for the selected day.");

            }, (error) => {
                this.setBanner("danger", "Fetching day log state failed! " + error);
            });
        },

        submitDayLog(event) {
            event.preventDefault();
            event.stopPropagation();

            this.setBanner();

            for (let i in this.logMetrics) {
                if (typeof (this.logMetrics[i]) === "string") {
                    this.logMetrics[i] = parseInt(this.logMetrics[i]);
                }
            }
            let reqBody = {
                "date": this.logDate + "T00:00:00Z",
                "metrics": this.logMetrics,
                "notes": this.logNotes
            };

            this.performDayLogRequest("/daylog/data/daylog", "POST", reqBody, () => {
                this.setBanner("success", "Day log submitted!");

            }, (error) => {
                this.setBanner("danger", "Day log submission failed! " + error);
            });
        },

        performDayLogRequest(url, method, body, successFunc, errorFunc) {
            axios({
                method: method,
                url: process.env.VUE_APP_API_HOST + url,
                data: JSON.stringify(body)
            })
                .then((resp) => {
                    if (successFunc) {
                        successFunc(resp.data);
                    }
                })
                .catch((error) => {
                    if (errorFunc) {
                        errorFunc(error);
                    }
                    console.error(error);
                });
        }
    }
};
</script>

<style>
</style>