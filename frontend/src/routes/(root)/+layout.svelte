<script>
    import { goto } from '$app/navigation';

    import * as utils from "$lib/js/utils.js"
</script>

{#await utils.isUserAuthenticated() then authenticated}
    {#if authenticated === true}
        <slot></slot>
    {:else if authenticated === false}
        {goto("/signin")}
    {/if}
{/await}