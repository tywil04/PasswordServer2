<script>
    import { goto } from '$app/navigation';

    import * as crypto from "$lib/js/crypto.js"
    import * as utils from "$lib/js/utils.js"

    import Button from "$lib/components/Button.svelte";
    import TextInput from "$lib/components/TextInput.svelte";

    let formData = {
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

        if (!formData.email.valid || formData.email.value == "") {
            return
        }

        if (!formData.password.valid || formData.password.value == "") {
            return
        }

        if (!formData.passwordConfirm.valid || formData.passwordConfirm.value == "") {
            return 
        }

        if (formData.password.value !== formData.passwordConfirm.value) {
            error = "Passwords are not the same."
            return
        }

        let hashedEmail = await crypto.hash(formData.email.value)
        let masterKey = await crypto.generateMasterKey(formData.password.value, hashedEmail) // Derive a key via pbkdf2 from the users password and email using
        let masterHash = await crypto.generateMasterHash(formData.password.value, masterKey) // Derive bits via pbkdf2 from the masterkey and the users password (this is used for server-side auth)

        let databaseKey = await crypto.generateDatabaseKey() // generate random AES-256-CBC key
        let [iv, encryptedDatabaseKey] = await crypto.protectDatabaseKey(masterKey, databaseKey) // encrypt the key with masterkey

        let response = await fetch("/api/v1/auth/signup", {
            method: "POST",
            body: JSON.stringify({
                Email: formData.email.value,
                MasterHash: Array.from(new Uint8Array(masterHash)),
                ProtectedDatabaseKey: {
                    Iv: Array.from(new Uint8Array(iv)),
                    Key: Array.from(new Uint8Array(encryptedDatabaseKey)),
                },
                Config: crypto.config,
            })
        })
        let jsonResponse = await response.json()

        let success = jsonResponse.UserId !== undefined // an error response would not contain "UserId" instead it would contain "Error"

        if (success) {
            goto("/")
        } else {
            // reset data
            formData = {
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

            for (let error of jsonResponse["Error(s)"] ) {
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
            <TextInput label="Email" bind:value={formData.email.value} bind:valid={formData.email.valid} required validation="email" invalidText="Invalid email address." name="email" autocomplete="email" grow/>
    
            <div class="spacer big"/>
    
            <TextInput label="Password" bind:value={formData.password.value} bind:valid={formData.password.valid} required visibilityButton validation="password" invalidText="Invalid password." name="password" autocomplete="new-password" grow/>
    
            <div class="spacer big"/>

            <TextInput label="Password Confirm" bind:value={formData.passwordConfirm.value} bind:valid={formData.passwordConfirm.valid} required visibilityButton validation="password" invalidText="Invalid password confirm." name="passwordConfirm" autocomplete="password" grow/>
    
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