<script>
    import { goto } from '$app/navigation';

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
        passwordConfirm: {
            value: "",
            valid: true,
        }
    }

    let error = ""

    async function signup() {
        error = ""

        if (!data.email.valid || data.email.value == "") {
            return
        }

        if (!data.password.valid || data.password.value == "") {
            return
        }

        if (!data.passwordConfirm.valid || data.passwordConfirm.value == "") {
            return 
        }

        if (data.password.value !== data.passwordConfirm.value) {
            error = "Passwords are not the same."
            return
        }

        let masterKey = await crypto.generateMasterKey(data.password.value, data.email.value) // Derive a key via pbkdf2 from the users password and email using
        let masterHash = await crypto.generateMasterHash(data.password.value, masterKey) // Derive bits via pbkdf2 from the masterkey and the users password (this is used for server-side auth)

        let databaseKey = await crypto.generateDatabaseKey() // generate random AES-256-CBC key
        let [iv, encryptedDatabaseKey] = await crypto.protectDatabaseKey(masterKey, databaseKey) // encrypt the key with masterkey

        let response = await fetch("/api/v1/auth/signup", {
            method: "POST",
            body: JSON.stringify({
                Email: data.email.value,
                MasterHash: new Uint8Array(masterHash),
                ProtectedDatabaseKey: {
                    Iv: new Uint8Array(iv),
                    Key: new Uint8Array(encryptedDatabaseKey)
                },
            })
        })
        let jsonResponse = await response.json()

        let success = jsonResponse.UserId !== undefined // an error response would not contain "UserId" instead it would contain "Error"

        if (success) {
            goto("/")
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
    <title>Sign up</title>
</svelte:head>

<main class="page">
    <div class="outer">
        <div class="inner">
            <h1 class="title">Sign up</h1>
            
            <pre class="disclaimer">To sign in, click the button labled 'Sign in'
to be redirected to the correct page.

To sign up, ender your credentials. 

Your password must have at least 2
lowercase letters, uppercase letters, 
numbers and special characters.

Allowed special characters:
<code class="specialCharacters">!#$%^'"`&*-=_+&gt;&lt;?;:()&#123;&#125;[].,@</code>

Your password must also have a length
of 12 characters or more.</pre>
        </div>
    
        <div class="spacer verticalDesktopHorizontalMobile big"/>
    
        <form on:submit|preventDefault={signup} class="inner">
            <TextInput label="Email" bind:value={data.email.value} bind:valid={data.email.valid} required validation="email" invalidText="Invalid email address." name="email" autocomplete="email" grow/>
    
            <div class="spacer big"/>
    
            <TextInput label="Password" bind:value={data.password.value} bind:valid={data.password.valid} required visibilityButton validation="password" invalidText="Invalid password." name="password" autocomplete="new-password" grow/>
    
            <div class="spacer big"/>

            <TextInput label="Password Confirm" bind:value={data.passwordConfirm.value} bind:valid={data.passwordConfirm.valid} required visibilityButton validation="password" invalidText="Invalid password confirm." name="passwordConfirm" autocomplete="password" grow/>
    
            <div class="spacer big"/>

            {#if error !== ""}
                <span class="validationFailed">{error}</span>
                <div class="spacer big"/>
            {/if}

            <div class="buttonGroup">
                <Button submit grow variant="accent">Sign up</Button>
    
                <div class="spacer vertical big"/>
    
                <Button grow link href="/auth/signin">Sign in</Button>
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
        border-radius: calc(var(--borderRadius) * 2);
    }

    .buttonGroup {
        display: flex;
        flex-direction: row;
    }

    .disclaimer {
        padding: 0;
        margin: 0;
        font-family: var(--fontFamily);
    }

    .specialCharacters {
        color: var(--darkGray1);
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