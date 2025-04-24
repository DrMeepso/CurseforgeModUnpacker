<script lang="ts">
    import FileSelect from "./FileSelect.svelte";
    import Download from "./Download.svelte";
    import { EventsOn, EventsEmit } from '../wailsjs/runtime/runtime'

    enum State {
        FileSelect,
        Downloading,
        Finished,
    }

    const allOutput: string[] = $state([]);
    allOutput.push("Starting...");


    EventsOn("output", (output: string) => {
        console.log("Output: ", output);
        allOutput.push(output);
    });

    let currentState: State = $state(State.FileSelect);
    EventsOn("stateChange", (state_index: number) => {
        console.log("State changed to: ", state_index);
        currentState = state_index in State ? state_index : State.FileSelect;

        console.log("Current state: ", currentState);
        console.log("FileSelect state: ", State.FileSelect);
        console.log("Downloading state: ", State.Downloading);
        console.log("Finished state: ", State.Finished);

    });
    

</script>

<main>
    <h3>Curseforge Modpack Unpacker</h3>
    <hr>

    {#if currentState == State.FileSelect}
        <FileSelect/>
    {:else if currentState == State.Downloading}
        <Download allOutput={allOutput} />
    {:else if currentState == State.Finished}
        <p>Finished!</p>
    {/if}

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

    p {
        font: sans-serif;
    }

    .directory-selector {
        display: flex;
        justify-content: center;
        align-items: center;
        margin-top: 2px;
    }

    .directory-selector input {
        width: 100%;
        padding: 0.5em;
        border: 1px solid #ccc;
        border-radius: 4px;
        font-size: 1em;
    }
    .directory-selector button {
        margin-left: 0.5em;
        padding: 0.5em 1em;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 1em;
    }

    #download-button {
        display: flex;
        flex-direction: column;
        gap: 2px;
        margin-top: 10px;
        font: sans-serif;
    }

    #download-button p {
        margin: 0;
        font-size: 0.8em;
        color: #888;
    }

    .small-text {
        margin: 0;
        font-size: 0.8em;
        color: #888;
        width: 100%;
        text-align: left;
    }

    #download-button button {
        padding: 0.5em 1em;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 1em;
    }

    #download-button button:disabled {
        background-color: #ccc;
        cursor: not-allowed;
    }
    #download-button button:hover:not(:disabled) {
        background-color: #0056b3;
    }
    #download-button button:active:not(:disabled) {
        background-color: #004494;
    }

    #extra-options {
        display: flex;
        flex-direction: column;
        gap: 2px;
        margin-top: 0;
        font: sans-serif;
    }
    .extra-option {
        margin-left: 10px;
        display: flex;
        align-items: center;
        gap: 5px;
    }
    .extra-option input {
        width: 15px;
        height: 15px;
        cursor: pointer;
    }
    .extra-option label {
        font-size: 0.8em;
        color: #888;
        cursor: pointer;
    }
    .extra-option label:hover {
        color: #555;
    }
    .extra-option input:checked + label {
        color: #007bff;

    }
    input[type="checkbox"]:checked {
        background-color: #007bff; /* Blue background */
        border-color: #007bff;
    }


</style>
