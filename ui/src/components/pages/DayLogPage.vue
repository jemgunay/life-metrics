<template>
    <div class="col-md-12">
        <div class="card">
            <h5 class="card-header">Day Log</h5>
            <div class="card-body">
                <form>
                    <div class="form-group">
                        <label for="general-mood-input">General Mood</label>
                        <input type="range" class="form-control" id="general-mood-input" min="0" max="10" v-model="dayLogData['general_mood']">
                    </div>

                    <div class="form-group">
                        <label for="work-mood-input">Work Mood</label>
                        <input type="range" class="form-control" id="work-mood-input" min="0" max="10" v-model="dayLogData['work_mood']">
                    </div>

                    <div class="form-group">
                        <label for="sleep-input">Sleep Quality</label>
                        <input type="range" class="form-control" id="sleep-input" min="0" max="10" v-model="dayLogData['sleep_quality']">
                    </div>

                    <div class="form-group">
                        <label for="diet-input">Diet Quality</label>
                        <input type="range" class="form-control" id="diet-input" min="0" max="10" v-model="dayLogData['diet_quality']">
                    </div>

                    <div class="form-group">
                        <label for="water-input">Water Intake</label>
                        <input type="range" class="form-control" id="water-input" min="0" max="10" v-model="dayLogData['water_intake']">
                    </div>

                    <div class="form-group">
                        <label for="exercise-input">Exercise</label>
                        <input type="range" class="form-control" id="exercise-input" min="0" max="10" v-model="dayLogData['exercise']">
                    </div>

                    <div class="form-group">
                        <label for="physical-input">Physical Health</label>
                        <input type="range" class="form-control" id="physical-input" min="0" max="10" v-model="dayLogData['physical_health']">
                    </div>

                    <div class="form-group">
                        <label for="notes-input">Notes</label>
                        <textarea class="form-control" id="notes-input" v-model="dayLogData['notes']"></textarea>
                    </div>

                    <button type="submit" class="btn btn-primary">Submit</button>
                </form>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "DayLogPage",
    data() {
        return {
            dayLogData: {}
        };
    }
};
</script>

<style>
</style>