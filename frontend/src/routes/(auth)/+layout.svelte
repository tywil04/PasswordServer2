<script>
    import { goto } from '$app/navigation';

    import * as utils from "$lib/js/utils.js"
</script>

{#await utils.isUserAuthenticated() then authenticated}
    {#if authenticated === true}
        {goto("/")}
    {:else if authenticated === false}
        <slot></slot>
    {/if}
{/await}