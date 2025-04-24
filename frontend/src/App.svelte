<script lang="ts">
    import FileSelect from "./FileSelect.svelte";
    import Download from "./Download.svelte";
    import { EventsOn, EventsEmit } from '../wailsjs/runtime/runtime'
    import { OpenDonationPage } from '../wailsjs/go/main/App'

    enum State {
        FileSelect,
        Downloading
    }

    const UnpackerOutput: string[] = $state([]);
    UnpackerOutput.push("Starting...");


    EventsOn("output", (output: string) => {
        UnpackerOutput.push(output);
    });

    let currentState: State = $state(State.FileSelect);
    EventsOn("stateChange", (state_index: number) => {
        console.log("State changed to: ", state_index);
        currentState = state_index in State ? state_index : State.FileSelect;
    });

</script>

<main>
    <h3>Curseforge Modpack Unpacker</h3>
    <hr>

    {#if currentState == State.FileSelect}
        <FileSelect/>
    {:else if currentState == State.Downloading}
        <Download allOutput={UnpackerOutput} />
    {/if}

    <p id="kofi-support">Support me on <a onclick={OpenDonationPage}>Kofi</a></p>

</main>


<style>

    @import url('https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&display=swap');
    
    main {
        text-align: center;
        padding: 1em;
        max-width: 360px;
        margin: 0 auto;
        font-family: "Montserrat", sans-serif;
        font-optical-sizing: auto;
        font-weight: 600;
        font-style: normal;
    }
    
    h3 {
        font-size: 1.5em;
        font-weight: normal;
        margin: 0;
        font-weight: 700;
    }

    #kofi-support
    {
        margin-top: 20px;
        font-size: 0.8em;
        color: #888;
        position: absolute;
        left: 5px;
        bottom: 5px;
        margin: 0;
    }

    #kofi-support a {
        color: #007bff;
        text-decoration: underline;
        cursor: pointer;
    }

</style>
