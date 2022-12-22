<script>
    import * as crypto from "$lib/js/crypto.js"
    import * as utils from "$lib/js/utils.js"

    import Button from "$lib/components/Button.svelte";
    import TextInput from "$lib/components/TextInput.svelte";

    let data = {
        email: {
            value: "",
            valid: false,
        },
        password: {
            value: "",
            valid: false,
        },
    }

    async function signin() {
        if (!data.email.valid) {
            console.log("Email not valid!")
            return
        }

        if (!data.password.valid) {
            console.log("Password not valid!")
            return
        }

        let masterKey = await crypto.generateMasterKey(data.password.value, data.email.value) // Derive a key via pbkdf2 from the users password and email using
        let masterHash = utils.arrayBufferToHex(await crypto.generateMasterHash(data.password.value, masterKey)) // Derive bits via pbkdf2 from the masterkey and the users password (this is used for server-side auth)

        let response = await fetch("/api/v1/auth/signin", {
                method: "POST",
                body: JSON.stringify({
                Email: data.email.value,
                MasterHash: masterHash,
            })
        })
        let jsonResponse = await response.json()

        // jsonResponse.Authenticated is only used as a quick way to see if a user is authenticated, authentication is used server-side, this value means nothing
        if (jsonResponse.Authenticated) {
            let protectedDatabaseKeyRequest = await fetch("/api/v1/user/protectedDatabaseKey", {method: "GET"})
            let protectedDatabaseKey = (await protectedDatabaseKeyRequest.json()).ProtectedDatabaseKey

            sessionStorage.setItem("ProtectedDatabaseKey", protectedDatabaseKey)
            sessionStorage.setItem("MasterKey", utils.arrayBufferToHex(await crypto.exportRawKey(masterKey)))

            window.location = "/"
        } else {
            window.location.reload()
        }
    }
</script>

<svelte:head>
    <title>Sign in</title>
</svelte:head>

<main>
    <div class="rowColumnContainer">
        <div class="container">
            <h1>Sign in</h1>
            
            <pre>To sign up, click the button labled 'Sign up'
to be redirected to the correct page.

To sign in, ender your credentials.</pre>
        </div>
    
        <div class="spacer verticalDesktopHorizontalMobile big"/>
    
        <form on:submit|preventDefault={signin} class="container">
            <label class="inputLabel" for="email">Email</label>
    
            <div class="spacer small"/>
    
            <TextInput bind:value={data.email.value} bind:valid={data.email.valid} validation="email" name="email" autofocus autocomplete="newemail" grow/>
    
            <div class="spacer big"/>
    
            <label class="inputLabel" for="password">Password</label>
    
            <div class="spacer small"/>
    
            <TextInput bind:value={data.password.value} bind:valid={data.password.valid} visibilityButton validation="password" name="password" autocomplete="newpassword" grow/>
    
            <div class="spacer big"/>

            <div class="buttonContainer">
                <Button submit grow variant="accent">Sign in</Button>
    
                <div class="spacer vertical big"/>
    
                <Button grow link href="/auth/signup">Sign up</Button>
            </div>
        </form>
    </div>
</main>

<style>
    main {
        width: 100%;
        height: 100%;
        background-color: var(--lightGray5);
        display: flex;
        flex-direction: column;
        justify-content: center;
    }

    .rowColumnContainer {
        display: flex;
        flex-direction: row;
        width: fit-content;
        height: fit-content;
        margin-left: auto;
        margin-right: auto;
    }

    .container > h1 {
        margin-top: 0;
        margin-bottom: 10px;
    }

    .container {
        width: fit-content;
        height: fit-content;
        padding: 20px;
        border: 1px solid var(--lightGray4);
        background-color: white;
        border-radius: var(--borderRadius);
    }

    label.inputLabel {
        text-align: left;
    }

    label.inputLabel::after {
        font-size: x-small;
        content: " (required)";
        color: var(--darkGray1);
    }

    div.spacer {
        padding: 0;
        margin: 0;
    }

    div.spacer.big:not(.vertical) {
        height: 15px;
    }

    div.spacer.small:not(.vertical) {
        height: 2px;
    }

    div.spacer.vertical.big {
        width: 10px;
    }

    div.spacer.verticalDesktopHorizontalMobile.big {
        width: 15px;
    }

    div.buttonContainer {
        display: flex;
        flex-direction: row;
    }

    .container > pre {
        padding: 0;
        margin: 0;
        font-family: var(--defaultFontFamily);
    }

    @media only screen and (max-width: 600px) {
        .rowColumnContainer {
            flex-direction: column;
        }

        .rowColumnContainer > * {
            flex-grow: 1;
            width: auto;
        }

        div.spacer.verticalDesktopHorizontalMobile.big {
            height: 15px;
        }
    }
</style>