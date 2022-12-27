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
    }

    let error = ""

    async function signin(event, optionalConfig=null, optionalOldConfig=null) {
        if (!formData.email.valid) {
            console.log("Email not valid!")
            return
        }

        if (!formData.password.valid) {
            console.log("Password not valid!")
            return
        }

        let config = optionalConfig === null ? crypto.config: optionalConfig

        let hashedEmail = await crypto.hash(formData.email.value, config)
        let masterKey = await crypto.generateMasterKey(formData.password.value, hashedEmail, config) // Derive a key via pbkdf2 from the users password and email using
        let masterHash = await crypto.generateMasterHash(formData.password.value, masterKey, config) // Derive bits via pbkdf2 from the masterkey and the users password (this is used for server-side auth)

        let response = await fetch("/api/v1/auth/signin", {
            method: "POST",
            body: JSON.stringify({
                Email: formData.email.value,
                MasterHash: Array.from(new Uint8Array(masterHash)),
                Config: config,
            })
        })
        let jsonResponse = await response.json()

        // jsonResponse.Authenticated is only used as a quick way to see if a user is authenticated, authentication is used server-side, this value means nothing
        if (jsonResponse.Authenticated) {
            let pdk = jsonResponse.ProtectedDatabaseKey
            let userId = jsonResponse.UserId
            sessionStorage.setItem("PasswordServer2:ProtectedDatabaseKey", JSON.stringify(pdk))
            sessionStorage.setItem("PasswordServer2:Email", hashedEmail)
            sessionStorage.setItem("PasswordServer2:UserId", userId)

            if (optionalConfig !== null) {
                let masterKey = await crypto.generateMasterKey(formData.password.value, hashedEmail, optionalOldConfig)
                let dk = crypto.unprotectDatabaseKey(masterKey, [pdk.Iv, pdk.Key], optionalOldConfig)
                let newDk = crypto.protectDatabaseKey(masterKey, dk, config)

                await fetch("/api/v1/users/configprofiles/new", {
                    method: "POST",
                    body: JSON.stringify({
                        UserId: userId,
                        MasterHash: Array.from(new Uint8Array(masterHash)),
                        ProtectedDatabaseKey: {
                            "Iv": Array.from(new Uint8Array(newDk[0])),
                            "Key": Array.from(new Uint8Array(newDk[1])),
                        },
                        Config: config,
                    })
                })

                dk = null
            }

            goto("/")
        } else if (!jsonResponse.Authenticated && jsonResponse.NewConfigRequired === true) {
            let oldConfig = jsonResponse.OldConfigs[0]
            console.log(oldConfig)
            signin(event, oldConfig, config)
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
            <TextInput label="Email" bind:value={formData.email.value} bind:valid={formData.email.valid} required validation="email" invalidText="Invalid email address." name="email" autocomplete="email" grow/>
    
            <div class="spacer big"/>
    
            <TextInput label="Password" bind:value={formData.password.value} bind:valid={formData.password.valid} required visibilityButton validation="password" invalidText="Invalid password." name="password" autocomplete="newpassword" grow/>
    
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