<script>
    export let variant = ""
    export let submit = false
    export let grow = false
    export let link = false
    export let href

    export let onClick = () => {}
</script>


{#if !link} 
    <button on:click={onClick} type={submit ? "submit": "button"} class={"trigger " + variant} class:grow={grow}>
        <slot></slot>
    </button>
{:else}
    <a href={href} class={"trigger " + variant} class:grow={grow}>
        <slot></slot>
    </a>
{/if}

<style>
    .trigger {
        background-color: var(--lightGray5);
        color: var(--textLight);
        border: 1px solid var(--lightGray4);
        border-radius: var(--borderRadius);
        padding: var(--defaultPadding);
        transition-duration: var(--animationDuration);
        text-decoration: none;
        text-align: center;
        line-height: 1.3em;
    }

    .accent {
        background-color: var(--blue);
        border: 1px solid var(--blue);
        color: var(--textDark);
    }

    .minimal {
        background-color: transparent;
        border: 1px solid transparent;
    }

    .dangerous {
        background-color: var(--red);
        border: 1px solid var(--red);
        color: var(--textDark);
    }

    .minimal:hover:not(:active) {
        background-color: var(--lightGray5);
        border: 1px solid var(--lightGray5)
    }

    .trigger:hover, .trigger:active {
        cursor: pointer;
    }

    .trigger:hover:not(:active) {
        opacity: 75%;;
    }

    .grow {
        flex-grow: 1;
    }
</style>