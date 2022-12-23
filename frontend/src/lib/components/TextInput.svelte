<script>
    export let value = null
    export let validation = ""
    export let invalidText = ""
    export let valid = true
    export let visibilityButton = false
    export let grow = false
    export let label = ""
    export let name = ""
    export let required = false 
    export let autocomplete = ""
    export let autocapitalize = false

    let input
    let forceVisible = false
    let validations = {
        "email": /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/,
        "password": /^((?=.*([A-Z]){2,})(?=.*([!#$%^'"`&*-=_+><?;:(){}\[\].,@]){2,})(?=.*([0-9]){2,})(?=.*([a-z]){2,})).{12,}$/ // requires 2 lowercase, 2 uppercase, 2 numbers, 2 special characters and at least 12 characters
    }
    let selectedValidation = validations[validation.toLowerCase()]

    let setType = (node) => {
        if (validation === "password" && forceVisible) {
            node.type = "text"
        } else if (validation.match(/^(||url|email|password)$/) && !forceVisible) { // allow validation as "", "url", "email", "password"
            node.type = validation === "" ? "text": validation // if validation is blank, set type to text otherwise set type to approved validation
        }
    }

    let onInput = () => {
        let regex = new RegExp(selectedValidation)
        valid = regex.test(value)
    }

    let toggleVisibility = (event) => {
        forceVisible = !forceVisible
        setType(input)
        event.target.innerText = forceVisible ? "Hide": "Show"
    }
</script>

<div class="outer">
    {#if label}
        <label class="inputLabel" for={name}>
            {label} {#if required}<span class="subtleText">(required)</span>{/if}
        </label>
    {/if}
    
    <div class="inner" class:visibilityButton={visibilityButton}>
        <input class="input" bind:this={input} bind:value={value} on:input={onInput} {required} {name} {autocomplete} {autocapitalize} class:grow={grow} use:setType/>
        {#if visibilityButton}
            <button class="button" type="button" on:click={toggleVisibility}>Show</button>
        {/if}
    </div>

    {#if !valid && invalidText != ""}
        <span class="validationFailedText">{invalidText}</span>
    {/if}
</div>

<style>
    .subtleText {
        font-size: x-small;
        color: var(--darkGray1);
    }

    .validationFailedText {
        margin-top: 2px;
        font-size: x-small;
        color: var(--red);
    }

    .inputLabel {
        margin-bottom: 2px;
    }

    .outer {
        display: flex;
        flex-direction: column  ;
        height: fit-content;
    }

    .inner {
        display: flex;
        flex-direction: row;
        height: fit-content;
    }

    .input {
        background-color: var(--lightGray5);
        border: 1px solid var(--lightGray4);
        border-radius: var(--borderRadius);
        padding: var(--defaultPadding);
        margin: 0;
    }

    .input.grow {
        flex-grow: 1;
    }

    .input:focus, .input:hover {
        outline: none;
    }

    .inner.visibilityButton .input {
        border-radius: var(--borderRadius) 0px 0px var(--borderRadius);
    }

    .inner.visibilityButton .button {
        background-color: var(--lightGray5);
        border: 1px solid var(--lightGray4);
        border-left: 0px;
        border-radius: 0px var(--borderRadius) var(--borderRadius) 0px;
        padding: var(--defaultPadding);
        font-family: monospace;
        font-size: x-small;
        color: var(--darkGray1);
    }

    .inner.visibilityButton .button:hover, .inner.visibilityButton .button:active {
        cursor: pointer;
    }

    .inner.visibilityButton .button:hover:not(:active) {
        opacity: 75%;;
    }
</style>