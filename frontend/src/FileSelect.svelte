<script lang="ts">
  import { onMount } from 'svelte';
    import { ZipFileDialog, VerifyModpackFile, ShowErrorDialog, FolderDialog, RunUnpack, SetIncludeOverrides } from '../wailsjs/go/main/App'
    
    let zipFilePath: string = $state("")
    let isValidModpack: boolean = $state(false)
    let modpackDisplayName: string = $state("")

    let extractPath: string = $state("")

    let includeOverrides: boolean = $state(false)

    $effect(() => {
        verifyZipFilePath(zipFilePath); // This will be called whenever zipFilePath changes
    });

    async function verifyZipFilePath(filePath: string) {
    
        if (filePath.length == 0) {
            isValidModpack = false;
            return;
        }

        let isValid = await VerifyModpackFile(filePath);
        if (isValid.split(":")[0].toLocaleLowerCase() == "error") {
            ShowErrorDialog(isValid)
            isValidModpack = false;
            return;
        }
        isValidModpack = true;
        modpackDisplayName = isValid.split(":")[1];
    }

    let IncludeOverides: HTMLInputElement | null = null;
    onMount(() => {
        IncludeOverides!.oninput = () => {
            includeOverrides = IncludeOverides?.checked || false;
            SetIncludeOverrides(includeOverrides);
        };
    });

</script>

<main>
    <p>Select modpack zip file</p>

    <p class="small-text">Curseforge modpack file to download (.zip)</p>
    <div class="directory-selector">
        <input type="text" bind:value={zipFilePath} placeholder="Modpack file" />
        <button onclick={async () => {
            const path = await ZipFileDialog()
            zipFilePath = path
        }}>...</button>
    </div>

    <br>

    <p class="small-text">Folder for the mods to be unpacked too</p>
    <div class="directory-selector">
        <input type="text" bind:value={extractPath} placeholder="Download folder" />
        <button onclick={async () => {
            const path = await FolderDialog()
            extractPath = path
        }}>...</button>
    </div>

    <br>
    <p class="small-text">Extra Options:</p>
    <div id="extra-options">
        <div class="extra-option">
            <input type="checkbox" id="include-overrides" bind:group={includeOverrides} bind:this={IncludeOverides} />
            <label for="include-overrides">Include overrides</label>
        </div>
    </div>

    <div id="download-button">
        {#if isValidModpack}
            <p>{modpackDisplayName}</p>
            {#if extractPath.length > 0}
                <button onclick={() => {
                    RunUnpack(zipFilePath, extractPath)
                }}>Download</button>
            {:else}
                <button disabled>Download</button>
            {/if}
        {:else}
            <p>Please select a modpack file</p>
            <button disabled>Download</button>
        {/if}
    </div>
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
