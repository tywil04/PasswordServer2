<slot></slot>

<style>
    :root {
        --red: rgba(255, 59, 48, 100%);
        --orange: rgba(255, 149, 0, 100%);
        --yellow: rgba(255, 204, 0, 100%);
        --green: rgba(52, 199, 89, 100%);
        --mint: rgba(0, 199, 190, 100%);
        --teal: rgba(48, 176, 199, 100%);
        --cyan: rgba(50, 173, 230, 100%);
        --blue: rgba(0, 122, 255, 100%);
        --indigo: rgba(88, 88, 214, 100%);
        --purple: rgba(175, 82, 222, 100%);
        --pink: rgba(255, 45, 85, 100%);

        --lightGray5: rgba(242, 242, 247, 100%);
        --lightGray4: rgba(229, 229, 234, 100%);
        --lightGray3: rgba(209, 209, 214, 100%);
        --lightGray2: rgba(199, 199, 204, 100%);
        --lightGray1: rgba(174, 174, 179, 100%);
        --gray: rgba(142, 142, 147, 100%);
        --darkGray1: rgba(99, 99, 102, 100%);
        --darkGray2: rgba(72, 72, 74, 100%);
        --darkGray3: rgba(58, 58, 60, 100%);
        --darkGray4: rgba(44, 44, 46, 100%);
        --darkGray5: rgba(28, 28, 30, 100%);

        --textLight: black;
        --textDark: white;

        --shadowVerySmall: 0 0 0 1px rgba(0, 0, 0, 0.05); 
        --shadowSmall: 0 1px 2px 0 rgba(0, 0, 0, 0.05); 
        --shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06); 
        --shadowMedium: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06); 
        --shadowLarge: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05); 
        --shadowVeryLarge: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04); 

        --animationDuration: 100ms;

        --borderRadius: 6px;
        --defaultPadding: 3px 7px;
        --defaultFontSize: small;
        --defaultFontFamiy: sans-serif;

        font-size: var(--defaultFontSize);
        font-family: var(--defaultFontFamiy);
    }
</style>