<script>
    import { processProps } from "$lib/js/utils.js"

    export let value = null
    export let validation = ""
    export let valid = true
    export let visibilityButton = false
    export let grow = false

    let props = processProps($$props, ["value", "validation", "valid", "visibilityButton"])
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

<div class:visibilityButton={visibilityButton}>
    <input bind:this={input} bind:value={value} on:input={onInput} {...props} use:setType class:grow={grow}/>
    {#if visibilityButton}
        <button type="button" on:click={toggleVisibility}>Show</button>
    {/if}
</div>

<style>
    div {
        display: flex;
        flex-direction: row;
        height: fit-content;
    }

    input {
        background-color: var(--lightGray5);
        border: 1px solid var(--lightGray4);
        border-radius: var(--borderRadius);
        padding: var(--defaultPadding);
    }

    input.grow {
        flex-grow: 1;
    }

    input:focus, input:hover {
        outline: none;
    }

    div.visibilityButton input {
        border-radius: var(--borderRadius) 0px 0px var(--borderRadius);
    }

    div.visibilityButton button {
        background-color: var(--lightGray5);
        border: 1px solid var(--lightGray4);
        border-left: 0px;
        border-radius: 0px var(--borderRadius) var(--borderRadius) 0px;
        padding: var(--defaultPadding);
        font-family: monospace;
        font-size: x-small;
        color: var(--darkGray1);
    }

    div.visibilityButton button:hover, div:nth-child(2) button:active {
        cursor: pointer;
    }

    div.visibilityButton button:hover:not(:active) {
        opacity: 75%;;
    }
</style>