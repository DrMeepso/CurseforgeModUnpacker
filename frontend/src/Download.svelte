<script lang="ts">
  import { onMount } from 'svelte';
    import { ZipFileDialog, VerifyModpackFile, ShowErrorDialog, FolderDialog } from '../wailsjs/go/main/App'
    import { EventsEmit, EventsOn } from '../wailsjs/runtime/runtime'
    let { allOutput } = $props();

    let progressValue: number = $state(0);
    
    $effect(() => {
        const length = allOutput.length;
        scrollToBottom(allOutput);
    })

    function scrollToBottom(trigger: any) {
        const outputHolder = document.getElementById("output-holder");
        if (outputHolder) {
            console.log("Scrolling to bottom...");
            outputHolder.scrollTop = outputHolder.scrollHeight;
        }
    }

    let progressBar: HTMLProgressElement | null = null;
    onMount(() => {
        EventsOn("progress", (progress: number) => {
            if (progressBar) {
                progressValue = progress;
                progressBar.value = progressValue;
            }
        });
    });

    let isFinished : boolean = $state(false);

    EventsOn("finished", () => {
        if (progressBar) {
            progressBar.value = 100;
            progressValue = 100;
        }
        // add a bit of delay to show the progress bar at 100%
        setTimeout(() => {
            isFinished = true;
        }, 100);
    });

</script>

<main>
    <p>Please wait, Unpacking...</p>

    <div id="output-holder">
        {#each allOutput as output}
            <p>{output}</p>
        {/each}
    </div>

    <progress value="0" max="100" style="width: 100%; margin-top: 10px;" bind:this={progressBar}></progress>
    <p>{progressValue}%</p>

    {#if isFinished}
        <button onclick={() => {
            EventsEmit("stateChange", 0); // Change to finished state
        }} style="margin-top: 10px;">Done</button>
    {:else}
        <button style="margin-top: 10px;" disabled>Done</button>
    {/if}

</main>


<style>

    @import url('https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&display=swap');
    
    main {
        text-align: center;
        max-width: 360px;
        margin: 0 auto;
        font-family: "Montserrat", sans-serif;
        font-optical-sizing: auto;
        font-weight: 600;
        font-style: normal;
    }

    #output-holder {
        display: flex;
        flex-direction: column;
        gap: 2px;
        margin-top: 10px;
        font: sans-serif;
        border: 1px solid #ccc;
        border-radius: 4px;
        padding: 5px;
        width: calc(100% - 10px);
        height: 200px;
        overflow-y: scroll;
    }

    #output-holder p {
        margin: 0;
        font-size: 0.8em;
        color: #888;
        text-align: left;
    }

    progress {
        width: 100%;
        height: 20px;
        margin-top: 10px;
    }

    progress::-webkit-progress-bar {
        background-color: #f3f3f3;
        border-radius: 4px;
    }
    progress::-webkit-progress-value {
        background-color: #007bff;
        border-radius: 4px;
    }

    p {
        margin: 0;
        font-size: 0.8em;
        color: #888;
    }

    button {
        padding: 0.5em 1em;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 1em;
        width: 100%;
    }

    button:disabled {
        background-color: #ccc;
        cursor: not-allowed;
    }
    button:hover:not(:disabled) {
        background-color: #0056b3;
    }
    button:active:not(:disabled) {
        background-color: #004494;
    }

</style>
