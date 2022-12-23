<script>
    import * as crypto from "$lib/js/crypto.js"
    import * as utils from "$lib/js/utils.js"

    import Button from "$lib/components/Button.svelte";
    import TextInput from "$lib/components/TextInput.svelte";

    let data = {
        email: {
            value: "",
            valid: true,
        },
        password: {
            value: "",
            valid: true,
        },
    }

    let error = ""

    async function signin() {
        if (!data.email.valid) {
            console.log("Email not valid!")
            return
        }

        if (!data.password.valid) {
            console.log("Password not valid!")
            return
        }

        let hashedEmail = await crypto.hash(data.email.value)
        let masterKey = await crypto.generateMasterKey(data.password.value, hashedEmail) // Derive a key via pbkdf2 from the users password and email using
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
            sessionStorage.setItem("Email", hashedEmail)

            window.location = "/"
        } else {
            // reset data
            data = {
                email: {
                    value: "",
                    valid: true,
                },
                password: {
                    value: "",
                    valid: true,
                },
                passwordConfirm: {
                    value: "",
                    valid: true,
                }
            }

            let tempError = []

            for (let error of jsonResponse["Errors(s)"]) {
                tempError.push(error["Message"])
            }

            error = tempError.join("\n")
        }
    }
</script>

<svelte:head>
    <title>Sign in</title>
</svelte:head>

<main class="page">
    <div class="outer">
        <div class="inner">
            <h1 class="title">Sign in</h1>
            
            <pre class="disclaimer">To sign up, click the button labled 'Sign up'
to be redirected to the correct page.

To sign in, ender your credentials.</pre>
        </div>
    
        <div class="spacer verticalDesktopHorizontalMobile big"/>
    
        <form on:submit|preventDefault={signin} class="inner">
            <TextInput label="Email" bind:value={data.email.value} bind:valid={data.email.valid} required validation="email" invalidText="Invalid email address." name="email" autocomplete="email" grow/>
    
            <div class="spacer big"/>
    
            <TextInput label="Password" bind:value={data.password.value} bind:valid={data.password.valid} required visibilityButton validation="password" invalidText="Invalid password." name="password" autocomplete="newpassword" grow/>
    
            <div class="spacer big"/>

            {#if error !== ""}
                <span class="validationFailed">{error}</span>
                <div class="spacer big"/>
            {/if}

            <div class="buttonGroup">
                <Button submit grow variant="accent">Sign in</Button>
    
                <div class="spacer vertical big"/>
    
                <Button grow link href="/auth/signup">Sign up</Button>
            </div>
        </form>
    </div>
</main>

<style>
    .page {
        width: 100%;
        height: 100%;
        background-color: var(--lightGray5);
        display: flex;
        flex-direction: column;
        justify-content: center;
    }

    .validationFailed {
        color: var(--red);
    }

    .outer {
        display: flex;
        flex-direction: row;
        width: fit-content;
        height: fit-content;
        margin-left: auto;
        margin-right: auto;
    }

    .title {
        margin-top: 0;
        margin-bottom: 10px;
    }

    .inner {
        width: fit-content;
        height: fit-content;
        padding: 20px;
        border: 1px solid var(--lightGray4);
        background-color: white;
        border-radius: var(--borderRadius);
    }

    .buttonGroup {
        display: flex;
        flex-direction: row;
    }

    .disclaimer {
        padding: 0;
        margin: 0;
        font-family: var(--defaultFontFamily);
    }

    @media only screen and (max-width: 600px) {
        .outer {
            flex-direction: column;
        }

        .outer > * {
            flex-grow: 1;
            width: auto;
        }
    }
</style>