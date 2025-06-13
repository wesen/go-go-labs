Title: The kitty remote control protocol

URL Source: https://sw.kovidgoyal.net/kitty/rc_protocol/

Markdown Content:
The kitty remote control protocol - kitty

*   [Terminal protocol extensions](https://sw.kovidgoyal.net/kitty/protocol-extensions/)- [x] Toggle navigation of Terminal protocol extensions  
    *   [Colored and styled underlines](https://sw.kovidgoyal.net/kitty/underlines/)
    *   [Terminal graphics protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/)- [x] Toggle navigation of Terminal graphics protocol  
        *   [Getting the window size](https://sw.kovidgoyal.net/kitty/graphics-protocol/#getting-the-window-size)
        *   [A minimal example](https://sw.kovidgoyal.net/kitty/graphics-protocol/#a-minimal-example)
        *   [The graphics escape code](https://sw.kovidgoyal.net/kitty/graphics-protocol/#the-graphics-escape-code)
        *   [Transferring pixel data](https://sw.kovidgoyal.net/kitty/graphics-protocol/#transferring-pixel-data)- [x] Toggle navigation of Transferring pixel data  
            *   [RGB and RGBA data](https://sw.kovidgoyal.net/kitty/graphics-protocol/#rgb-and-rgba-data)
            *   [PNG data](https://sw.kovidgoyal.net/kitty/graphics-protocol/#png-data)
            *   [Compression](https://sw.kovidgoyal.net/kitty/graphics-protocol/#compression)
            *   [The transmission medium](https://sw.kovidgoyal.net/kitty/graphics-protocol/#the-transmission-medium)- [x] Toggle navigation of The transmission medium  
                *   [Local client](https://sw.kovidgoyal.net/kitty/graphics-protocol/#local-client)
                *   [Remote client](https://sw.kovidgoyal.net/kitty/graphics-protocol/#remote-client)

            *   [Querying support and available transmission mediums](https://sw.kovidgoyal.net/kitty/graphics-protocol/#querying-support-and-available-transmission-mediums)

        *   [Display images on screen](https://sw.kovidgoyal.net/kitty/graphics-protocol/#display-images-on-screen)- [x] Toggle navigation of Display images on screen  
            *   [Controlling displayed image layout](https://sw.kovidgoyal.net/kitty/graphics-protocol/#controlling-displayed-image-layout)
            *   [Unicode placeholders](https://sw.kovidgoyal.net/kitty/graphics-protocol/#unicode-placeholders)
            *   [Relative placements](https://sw.kovidgoyal.net/kitty/graphics-protocol/#relative-placements)

        *   [Deleting images](https://sw.kovidgoyal.net/kitty/graphics-protocol/#deleting-images)
        *   [Suppressing responses from the terminal](https://sw.kovidgoyal.net/kitty/graphics-protocol/#suppressing-responses-from-the-terminal)
        *   [Requesting image ids from the terminal](https://sw.kovidgoyal.net/kitty/graphics-protocol/#requesting-image-ids-from-the-terminal)
        *   [Animation](https://sw.kovidgoyal.net/kitty/graphics-protocol/#animation)- [x] Toggle navigation of Animation  
            *   [Transferring animation frame data](https://sw.kovidgoyal.net/kitty/graphics-protocol/#transferring-animation-frame-data)
            *   [Controlling animations](https://sw.kovidgoyal.net/kitty/graphics-protocol/#controlling-animations)
            *   [Composing animation frames](https://sw.kovidgoyal.net/kitty/graphics-protocol/#composing-animation-frames)

        *   [Image persistence and storage quotas](https://sw.kovidgoyal.net/kitty/graphics-protocol/#image-persistence-and-storage-quotas)
        *   [Control data reference](https://sw.kovidgoyal.net/kitty/graphics-protocol/#control-data-reference)
        *   [Interaction with other terminal actions](https://sw.kovidgoyal.net/kitty/graphics-protocol/#interaction-with-other-terminal-actions)

    *   [Comprehensive keyboard handling in terminals](https://sw.kovidgoyal.net/kitty/keyboard-protocol/)- [x] Toggle navigation of Comprehensive keyboard handling in terminals  
        *   [Quickstart](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#quickstart)
        *   [An overview](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#an-overview)- [x] Toggle navigation of An overview  
            *   [Key codes](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#key-codes)
            *   [Modifiers](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#modifiers)
            *   [Event types](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#event-types)
            *   [Text as code points](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#text-as-code-points)
            *   [Non-Unicode keys](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#non-unicode-keys)

        *   [Progressive enhancement](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#progressive-enhancement)- [x] Toggle navigation of Progressive enhancement  
            *   [Disambiguate escape codes](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#disambiguate-escape-codes)
            *   [Report event types](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#report-event-types)
            *   [Report alternate keys](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#report-alternate-keys)
            *   [Report all keys as escape codes](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#report-all-keys-as-escape-codes)
            *   [Report associated text](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#report-associated-text)

        *   [Detection of support for this protocol](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#detection-of-support-for-this-protocol)
        *   [Legacy key event encoding](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#legacy-key-event-encoding)- [x] Toggle navigation of Legacy key event encoding  
            *   [Legacy functional keys](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#legacy-functional-keys)
            *   [Legacy text keys](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#legacy-text-keys)

        *   [Functional key definitions](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#functional-key-definitions)
        *   [Legacy ctrl mapping of ASCII keys](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#legacy-ctrl-mapping-of-ascii-keys)
        *   [Bugs in fixterms](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#bugs-in-fixterms)
        *   [Why xterm’s modifyOtherKeys should not be used](https://sw.kovidgoyal.net/kitty/keyboard-protocol/#why-xterm-s-modifyotherkeys-should-not-be-used)

    *   [The text sizing protocol](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/)- [x] Toggle navigation of The text sizing protocol  
        *   [Quickstart](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#quickstart)
        *   [The escape code](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#the-escape-code)
        *   [How it works](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#how-it-works)- [x] Toggle navigation of How it works  
            *   [Fractional scaling](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#fractional-scaling)

        *   [Fixing the character width issue for the terminal ecosystem](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#fixing-the-character-width-issue-for-the-terminal-ecosystem)
        *   [Wrapping and overwriting behavior](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#wrapping-and-overwriting-behavior)
        *   [Detecting if the terminal supports this protocol](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#detecting-if-the-terminal-supports-this-protocol)
        *   [Interaction with other terminal controls](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#interaction-with-other-terminal-controls)- [x] Toggle navigation of Interaction with other terminal controls  
            *   [Cursor movement](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#cursor-movement)
            *   [Editing controls](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#editing-controls)

        *   [The algorithm for splitting text into cells](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#the-algorithm-for-splitting-text-into-cells)- [x] Toggle navigation of The algorithm for splitting text into cells  
            *   [Unicode variation selectors](https://sw.kovidgoyal.net/kitty/text-sizing-protocol/#unicode-variation-selectors)

    *   [File transfer over the TTY](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/)- [x] Toggle navigation of File transfer over the TTY  
        *   [Overall design](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#overall-design)- [x] Toggle navigation of Overall design  
            *   [Sending files to the computer running the terminal emulator](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#sending-files-to-the-computer-running-the-terminal-emulator)
            *   [Receiving files from the computer running terminal emulator](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#receiving-files-from-the-computer-running-terminal-emulator)

        *   [Canceling a session](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#canceling-a-session)
        *   [Quieting responses from the terminal](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#quieting-responses-from-the-terminal)
        *   [File metadata](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#file-metadata)
        *   [Symbolic and hard links](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#symbolic-and-hard-links)- [x] Toggle navigation of Symbolic and hard links  
            *   [Sending links to the terminal emulator](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#sending-links-to-the-terminal-emulator)
            *   [Receiving links from the terminal emulator](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#receiving-links-from-the-terminal-emulator)

        *   [Transmitting binary deltas](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#transmitting-binary-deltas)- [x] Toggle navigation of Transmitting binary deltas  
            *   [Sending to the terminal emulator](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#sending-to-the-terminal-emulator)
            *   [Receiving from the terminal emulator](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#receiving-from-the-terminal-emulator)
            *   [The format of signatures and deltas](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#the-format-of-signatures-and-deltas)

        *   [Compression](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#compression)
        *   [Bypassing explicit user authorization](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#bypassing-explicit-user-authorization)
        *   [Encoding of transfer commands as escape codes](https://sw.kovidgoyal.net/kitty/file-transfer-protocol/#encoding-of-transfer-commands-as-escape-codes)

    *   [Desktop notifications](https://sw.kovidgoyal.net/kitty/desktop-notifications/)- [x] Toggle navigation of Desktop notifications  
        *   [Allowing users to filter notifications](https://sw.kovidgoyal.net/kitty/desktop-notifications/#allowing-users-to-filter-notifications)
        *   [Being informed when user activates the notification](https://sw.kovidgoyal.net/kitty/desktop-notifications/#being-informed-when-user-activates-the-notification)
        *   [Being informed when a notification is closed](https://sw.kovidgoyal.net/kitty/desktop-notifications/#being-informed-when-a-notification-is-closed)
        *   [Updating or closing an existing notification](https://sw.kovidgoyal.net/kitty/desktop-notifications/#updating-or-closing-an-existing-notification)
        *   [Automatically expiring notifications](https://sw.kovidgoyal.net/kitty/desktop-notifications/#automatically-expiring-notifications)
        *   [Adding icons to notifications](https://sw.kovidgoyal.net/kitty/desktop-notifications/#adding-icons-to-notifications)- [x] Toggle navigation of Adding icons to notifications  
            *   [Adding icons by transmitting icon data](https://sw.kovidgoyal.net/kitty/desktop-notifications/#adding-icons-by-transmitting-icon-data)

        *   [Adding buttons to the notification](https://sw.kovidgoyal.net/kitty/desktop-notifications/#adding-buttons-to-the-notification)
        *   [Playing a sound with notifications](https://sw.kovidgoyal.net/kitty/desktop-notifications/#playing-a-sound-with-notifications)
        *   [Querying for support](https://sw.kovidgoyal.net/kitty/desktop-notifications/#querying-for-support)
        *   [Specification of all keys used in the protocol](https://sw.kovidgoyal.net/kitty/desktop-notifications/#specification-of-all-keys-used-in-the-protocol)
        *   [Base64](https://sw.kovidgoyal.net/kitty/desktop-notifications/#base64)
        *   [Escape code safe UTF-8](https://sw.kovidgoyal.net/kitty/desktop-notifications/#escape-code-safe-utf-8)
        *   [Identifier](https://sw.kovidgoyal.net/kitty/desktop-notifications/#identifier)

    *   [Mouse pointer shapes](https://sw.kovidgoyal.net/kitty/pointer-shapes/)- [x] Toggle navigation of Mouse pointer shapes  
        *   [Setting the pointer shape](https://sw.kovidgoyal.net/kitty/pointer-shapes/#setting-the-pointer-shape)
        *   [Pushing and popping shapes onto the stack](https://sw.kovidgoyal.net/kitty/pointer-shapes/#pushing-and-popping-shapes-onto-the-stack)
        *   [Querying support](https://sw.kovidgoyal.net/kitty/pointer-shapes/#querying-support)
        *   [Interaction with other terminal features](https://sw.kovidgoyal.net/kitty/pointer-shapes/#interaction-with-other-terminal-features)
        *   [Pointer shape names](https://sw.kovidgoyal.net/kitty/pointer-shapes/#pointer-shape-names)
        *   [Legacy xterm compatibility](https://sw.kovidgoyal.net/kitty/pointer-shapes/#legacy-xterm-compatibility)

    *   [Unscrolling the screen](https://sw.kovidgoyal.net/kitty/unscroll/)
    *   [Color control](https://sw.kovidgoyal.net/kitty/color-stack/)- [x] Toggle navigation of Color control  
        *   [Saving and restoring colors](https://sw.kovidgoyal.net/kitty/color-stack/#saving-and-restoring-colors)
        *   [Setting and querying colors](https://sw.kovidgoyal.net/kitty/color-stack/#setting-and-querying-colors)- [x] Toggle navigation of Setting and querying colors  
            *   [Querying current color values](https://sw.kovidgoyal.net/kitty/color-stack/#querying-current-color-values)
            *   [Setting color values](https://sw.kovidgoyal.net/kitty/color-stack/#setting-color-values)
            *   [Color value encoding](https://sw.kovidgoyal.net/kitty/color-stack/#color-value-encoding)

    *   [Setting text styles/colors in arbitrary regions of the screen](https://sw.kovidgoyal.net/kitty/deccara/)
    *   [Copying all data types to the clipboard](https://sw.kovidgoyal.net/kitty/clipboard/)- [x] Toggle navigation of Copying all data types to the clipboard  
        *   [Reading data from the system clipboard](https://sw.kovidgoyal.net/kitty/clipboard/#reading-data-from-the-system-clipboard)
        *   [Writing data to the system clipboard](https://sw.kovidgoyal.net/kitty/clipboard/#writing-data-to-the-system-clipboard)
        *   [Support for terminal multiplexers](https://sw.kovidgoyal.net/kitty/clipboard/#support-for-terminal-multiplexers)

    *   [Miscellaneous protocol extensions](https://sw.kovidgoyal.net/kitty/misc-protocol/)- [x] Toggle navigation of Miscellaneous protocol extensions  
        *   [Simple save/restore of all terminal modes](https://sw.kovidgoyal.net/kitty/misc-protocol/#simple-save-restore-of-all-terminal-modes)
        *   [Independent control of bold and faint SGR properties](https://sw.kovidgoyal.net/kitty/misc-protocol/#independent-control-of-bold-and-faint-sgr-properties)
        *   [kitty specific private escape codes](https://sw.kovidgoyal.net/kitty/misc-protocol/#kitty-specific-private-escape-codes)

*   [Press mentions of kitty](https://sw.kovidgoyal.net/kitty/press-mentions/)- [x] Toggle navigation of Press mentions of kitty  
    *   [Video reviews](https://sw.kovidgoyal.net/kitty/press-mentions/#video-reviews)

[Back to top](https://sw.kovidgoyal.net/kitty/rc_protocol/#)

Toggle Light / Dark / Auto color theme

Toggle table of contents sidebar

 

The kitty remote control protocol[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#the-kitty-remote-control-protocol "Link to this heading")
===========================================================================================================================================

The kitty remote control protocol is a simple protocol that involves sending data to kitty in the form of JSON. Any individual command of kitty has the form:

<ESC>P@kitty-cmd<JSON object><ESC>\

Where `<ESC>` is the byte `0x1b`. The JSON object has the form:

{
 "cmd": "command name",
 "version": "<kitty version>",
 "no_response": "<Optional Boolean>",
 "kitty_window_id": "<Optional value of the KITTY_WINDOW_ID env var>",
 "payload": "<Optional JSON object>"
}

The `version` above is an array of the form `[0, 14, 2]`. If you are developing a standalone client, use the kitty version that you are developing against. Using a version greater than the version of the kitty instance you are talking to, will cause a failure.

Set `no_response` to `true` if you don’t want a response from kitty.

The optional payload is a JSON object that is specific to the actual command being sent. The fields in the object for every command are documented below.

As a quick example showing how easy to use this protocol is, we will implement the `@ ls` command from the shell using only shell tools.

First, run kitty as:

kitty -o allow_remote_control=socket-only --listen-on unix:/tmp/test

Now, in a different terminal, you can get the pretty printed `@ ls` output with the following command line:

echo -en '\eP@kitty-cmd{"cmd":"ls","version":[0,14,2]}\e\\' | socat - unix:/tmp/test | awk '{ print substr($0, 13, length($0) - 14) }' | jq -c '.data | fromjson' | jq .

There is also the statically compiled stand-alone executable `kitten` that can be used for this, available from the [kitty releases](https://github.com/kovidgoyal/kitty/releases) page:

kitten @ --help

Encrypted communication[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#encrypted-communication "Link to this heading")
-----------------------------------------------------------------------------------------------------------------------

Added in version 0.26.0.

When using the [`remote_control_password`](https://sw.kovidgoyal.net/kitty/conf/#opt-kitty.remote_control_password) option communication to the terminal is encrypted to keep the password secure. A public key is used from the [`KITTY_PUBLIC_KEY`](https://sw.kovidgoyal.net/kitty/glossary/#envvar-KITTY_PUBLIC_KEY) environment variable. Currently, only one encryption protocol is supported. The protocol number is present in [`KITTY_PUBLIC_KEY`](https://sw.kovidgoyal.net/kitty/glossary/#envvar-KITTY_PUBLIC_KEY) as `1`. The key data in this environment variable is [**Base-85**](https://datatracker.ietf.org/doc/html/rfc1924.html) encoded. The algorithm used is [Elliptic Curve Diffie Helman](https://en.wikipedia.org/wiki/Elliptic-curve_Diffie%E2%80%93Hellman) with the [X25519 curve](https://en.wikipedia.org/wiki/Curve25519). A time based nonce is used to minimise replay attacks. The original JSON command has the fields: `password` and `timestamp` added. The timestamp is the number of nanoseconds since the epoch, excluding leap seconds. Commands with a timestamp more than 5 minutes from the current time are rejected. The command is then encrypted using AES-256-GCM in authenticated encryption mode, with a symmetric key that is derived from the ECDH key-pair by running the shared secret through SHA-256 hashing, once. An IV of at least 96 bits of CSPRNG data is used. The tag for authenticated encryption **must** be at least 128 bits long. The tag **must** authenticate only the value of the `encrypted` field. A new command is created and transmitted that contains the fields:

{
 "version": "<kitty version>",
 "iv": "base85 encoded IV",
 "tag": "base85 encoded AEAD tag",
 "pubkey": "base85 encoded ECDH public key of sender",
 "encrypted": "The original command encrypted and base85 encoded"
}

Async and streaming requests[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#async-and-streaming-requests "Link to this heading")
---------------------------------------------------------------------------------------------------------------------------------

Some remote control commands require asynchronous communication, that is, the response from the terminal can happen after an arbitrary amount of time. For example, the `select-window` command requires the user to select a window before a response can be sent. Such command must set the field `async` in the JSON block above to a random string that serves as a unique id. The client can cancel an async request in flight by adding the `cancel_async` field to the JSON block. A async response remains in flight until the terminal sends a response to the request. Note that cancellation requests dont need to be encrypted as users must not be prompted for these and the worst a malicious cancellation request can do is prevent another sync request from getting a response.

Similar to async requests are _streaming_ requests. In these the client has to send a large amount of data to the terminal and so the request is split into chunks. In every chunk the JSON block must contain the field `stream` set to `true` and `stream_id` set to a random long string, that should be the same for all chunks in a request. End of data is indicated by sending a chunk with no data.

action[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#action "Link to this heading")
-------------------------------------------------------------------------------------

Fields are:

`action (required)`
The action to perform. Of the form: action [optional args…]

`match_window (optional)`
Window to run the action on

`self (default: False)`
Whether to use the window this command is run in as the active window

close-tab[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#close-tab "Link to this heading")
-------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
Which tab to close

`self (default: False)`
Boolean indicating whether to close the tab of the window the command is run in

`ignore_no_match (default: False)`
Boolean indicating whether no matches should be ignored or return an error

close-window[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#close-window "Link to this heading")
-------------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
Which window to close

`self (default: False)`
Boolean indicating whether to close the window the command is run in

`ignore_no_match (default: False)`
Boolean indicating whether no matches should be ignored or return an error

create-marker[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#create-marker "Link to this heading")
---------------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
Which window to create the marker in

`self (default: False)`
Boolean indicating whether to create marker in the window the command is run in

`marker_spec (optional)`
A list or arguments that define the marker specification, for example: [‘text’, ‘1’, ‘ERROR’]

detach-tab[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#detach-tab "Link to this heading")
---------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
Which tab to detach

`target_tab (default: None)`
Which tab to move the detached tab to the OS window it is run in

`self (default: False)`
Boolean indicating whether to detach the tab the command is run in

detach-window[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#detach-window "Link to this heading")
---------------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
Which window to detach

`target_tab (default: None)`
Which tab to move the detached window to

`self (default: False)`
Boolean indicating whether to detach the window the command is run in

`stay_in_tab (default: False)`
Boolean indicating focus should remain in the active tab after windows are moved

disable-ligatures[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#disable-ligatures "Link to this heading")
-----------------------------------------------------------------------------------------------------------

Fields are:

`strategy (required)`
One of `never`, `always` or `cursor`

`match_window (optional)`
Window to change opacity in

`match_tab (default: None)`
Tab to change opacity in

`all (default: False)`
Boolean indicating operate on all windows

env[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#env "Link to this heading")
-------------------------------------------------------------------------------

Fields are:

`env (required)`
Dictionary of environment variables to values. When a env var ends with = it is removed from the environment.

focus-tab[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#focus-tab "Link to this heading")
-------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
The tab to focus

focus-window[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#focus-window "Link to this heading")
-------------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
The window to focus

get-colors[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#get-colors "Link to this heading")
---------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
The window to get the colors for

`configured (default: False)`
Boolean indicating whether to get configured or current colors

get-text[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#get-text "Link to this heading")
-----------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
The window to get text from

`extent (default: screen)`
One of `screen`, `first_cmd_output_on_screen`, `last_cmd_output`, `last_visited_cmd_output`, `all`, or `selection`

`ansi (default: False)`
Boolean, if True send ANSI formatting codes

`cursor (optional)`
Boolean, if True send cursor position/style as ANSI codes

`wrap_markers (optional)`
Boolean, if True add wrap markers to output

`clear_selection (default: False)`
Boolean, if True clear the selection in the matched window

`self (default: False)`
Boolean, if True use window the command was run in

goto-layout[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#goto-layout "Link to this heading")
-----------------------------------------------------------------------------------------------

Fields are:

`layout (required)`
The new layout name

`match (default: None)`
Which tab to change the layout of

kitten[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#kitten "Link to this heading")
-------------------------------------------------------------------------------------

Fields are:

`kitten (required)`
The name of the kitten to run

`args (optional)`
Arguments to pass to the kitten as a list

`match (default: None)`
The window to run the kitten over

last-used-layout[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#last-used-layout "Link to this heading")
---------------------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
Which tab to change the layout of

`all (default: False)`
Boolean to match all tabs

launch[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#launch "Link to this heading")
-------------------------------------------------------------------------------------

Fields are:

`args (required)`
The command line to run in the new window, as a list, use an empty list to run the default shell

`match (default: None)`
The tab to open the new window in

`next_to (default: None)`
The window next to which to create the new window or empty string to use active window

`source_window (default: None)`
The window to use as source for data or empty string to use active window

`window_title (default: None)`
Title for the new window

`cwd (default: None)`
Working directory for the new window

`env (default: [])`
List of environment variables of the form NAME=VALUE

`var (default: [])`
List of user variables of the form NAME=VALUE

`os_panel (default: [])`
List of panel settings

`tab_title (default: None)`
Title for the new tab

`type (default: window)`
The type of window to open

`keep_focus (default: False)`
Boolean indicating whether the current window should retain focus or not

`copy_colors (default: False)`
Boolean indicating whether to copy the colors from the current window

`copy_cmdline (default: False)`
Boolean indicating whether to copy the cmdline from the current window

`copy_env (default: False)`
List of strings representing the local env vars

`hold (default: False)`
Boolean indicating whether to keep window open after cmd exits

`location (default: default)`
Where in the tab to open the new window

`allow_remote_control (default: False)`
Boolean indicating whether to allow remote control from the new window

`remote_control_password (default: [])`
A list of remote control passwords

`stdin_source (default: none)`
Where to get stdin for the process from

`stdin_add_formatting (default: False)`
Boolean indicating whether to add formatting codes to stdin

`stdin_add_line_wrap_markers (default: False)`
Boolean indicating whether to add line wrap markers to stdin

`spacing (default: [])`
A list of spacing specifications, see the docs for the set-spacing command

`marker (default: None)`
Specification for marker for new window, for example: “text 1 ERROR”

`logo (default: None)`
Path to window logo

`logo_position (default: None)`
Window logo position as string or empty string to use default

`logo_alpha (default: -1.0)`
Window logo alpha or -1 to use default

`self (default: False)`
Boolean, if True use tab the command was run in

`os_window_title (default: None)`
Title for OS Window

`os_window_name (default: None)`
WM_NAME for OS Window

`os_window_class (default: None)`
WM_CLASS for OS Window

`os_window_state (default: normal)`
The initial state for OS Window

`color (default: [])`
list of color specifications such as foreground=red

`watcher (default: [])`
list of paths to watcher files

`bias (default: 0.0)`
The bias with which to create the new window in the current layout

`wait_for_child_to_exit (default: False)`
Boolean indicating whether to wait and return child exit code

load-config[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#load-config "Link to this heading")
-----------------------------------------------------------------------------------------------

Fields are:

`paths (optional)`
List of config file paths to load

`override (default: [])`
List of individual config overrides

`ignore_overrides (default: False)`
Whether to apply previous overrides

ls[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#ls "Link to this heading")
-----------------------------------------------------------------------------

Fields are:

`all_env_vars (default: False)`
Whether to send all environment variables for every window rather than just differing ones

`match (default: None)`
Window to change colors in

`match_tab (default: None)`
Tab to change colors in

`self (default: False)`
Boolean indicating whether to list only the window the command is run in

new-window[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#new-window "Link to this heading")
---------------------------------------------------------------------------------------------

Fields are:

`args (required)`
The command line to run in the new window, as a list, use an empty list to run the default shell

`match (default: None)`
The tab to open the new window in

`title (default: None)`
Title for the new window

`cwd (default: None)`
Working directory for the new window

`keep_focus (default: False)`
Boolean indicating whether the current window should retain focus or not

`window_type (default: kitty)`
One of `kitty` or `os`

`new_tab (default: False)`
Boolean indicating whether to open a new tab

`tab_title (default: None)`
Title for the new tab

remove-marker[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#remove-marker "Link to this heading")
---------------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
Which window to remove the marker from

`self (default: False)`
Boolean indicating whether to detach the window the command is run in

resize-os-window[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#resize-os-window "Link to this heading")
---------------------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
Which window to resize

`self (default: False)`
Boolean indicating whether to close the window the command is run in

`incremental (default: False)`
Boolean indicating whether to adjust the size incrementally

`action (default: resize)`
The action to perform

`unit (default: cells)`
One of `cells` or `pixels`

`width (default: 0)`
Integer indicating desired window width

`height (default: 0)`
Integer indicating desired window height

`os_panel (optional)`
Settings for modifying the OS Panel

resize-window[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#resize-window "Link to this heading")
---------------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
Which window to resize

`self (default: False)`
Boolean indicating whether to resize the window the command is run in

`increment (default: 2)`
Integer specifying the resize increment

`axis (default: horizontal)`
One of `horizontal, vertical` or `reset`

run[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#run "Link to this heading")
-------------------------------------------------------------------------------

Fields are:

`data (required)`
Chunk of STDIN data, base64 encoded no more than 4096 bytes. Must send an empty chunk to indicate end of data.

`cmdline (required)`
The command line to run

`env (default: [])`
List of environment variables of the form NAME=VALUE

`allow_remote_control (default: False)`
A boolean indicating whether to allow remote control

`remote_control_password (default: [])`
A list of remote control passwords

scroll-window[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#scroll-window "Link to this heading")
---------------------------------------------------------------------------------------------------

> for unscrolling by lines, or ‘r’ for scrolling ot prompt.

Fields are:

`amount (required)`
The amount to scroll, a two item list with the first item being either a number or the keywords, start and end. And the second item being either ‘p’ for pages or ‘l’ for lines or ‘u’

`match (default: None)`
The window to scroll

select-window[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#select-window "Link to this heading")
---------------------------------------------------------------------------------------------------

Fields are:

`match (default: None)`
The tab to open the new window in

`self (default: False)`
Boolean, if True use tab the command was run in

`title (default: None)`
A title for this selection

`exclude_active (default: False)`
Exclude the currently active window from the list to pick

`reactivate_prev_tab (default: False)`
Reactivate the previously activated tab when finished

send-key[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#send-key "Link to this heading")
-----------------------------------------------------------------------------------------

Fields are:

`keys (required)`
The keys to send

`match (default: None)`
A string indicating the window to send text to

`match_tab (default: None)`
A string indicating the tab to send text to

`all (default: False)`
A boolean indicating all windows should be matched.

`exclude_active (default: False)`
A boolean that prevents sending text to the active window

send-text[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#send-text "Link to this heading")
-------------------------------------------------------------------------------------------

Fields are:

`data (required)`
The data being sent. Can be either: text: followed by text or base64: followed by standard base64 encoded bytes

`match (default: None)`
A string indicating the window to send text to

`match_tab (default: None)`
A string indicating the tab to send text to

`all (default: False)`
A boolean indicating all windows should be matched.

`exclude_active (default: False)`
A boolean that prevents sending text to the active window

`session_id (optional)`
A string that identifies a “broadcast session”

`bracketed_paste (default: disable)`
Whether to wrap the text in bracketed paste escape codes

set-background-image[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-background-image "Link to this heading")
-----------------------------------------------------------------------------------------------------------------

Fields are:

`data (required)`
Chunk of at most 512 bytes of PNG data, base64 encoded. Must send an empty chunk to indicate end of image. Or the special value - to indicate image must be removed.

`match (default: None)`
Window to change opacity in

`layout (default: configured)`
The image layout

`all (default: False)`
Boolean indicating operate on all windows

`configured (default: False)`
Boolean indicating if the configured value should be changed

set-background-opacity[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-background-opacity "Link to this heading")
---------------------------------------------------------------------------------------------------------------------

Fields are:

`opacity (required)`
A number between 0 and 1

`match_window (optional)`
Window to change opacity in

`match_tab (default: None)`
Tab to change opacity in

`all (default: False)`
Boolean indicating operate on all windows

`toggle (default: False)`
Boolean indicating if opacity should be toggled between the default and the specified value

set-colors[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-colors "Link to this heading")
---------------------------------------------------------------------------------------------

Fields are:

`colors (required)`
An object mapping names to colors as 24-bit RGB integers or null for nullable colors. Or a string for transparent_background_colors.

`match_window (optional)`
Window to change colors in

`match_tab (default: None)`
Tab to change colors in

`all (default: False)`
Boolean indicating change colors everywhere or not

`configured (default: False)`
Boolean indicating whether to change the configured colors. Must be True if reset is True

`reset (default: False)`
Boolean indicating colors should be reset to startup values

set-enabled-layouts[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-enabled-layouts "Link to this heading")
---------------------------------------------------------------------------------------------------------------

Fields are:

`layouts (required)`
The list of layout names

`match (default: None)`
Which tab to change the layout of

`configured (default: False)`
Boolean indicating whether to change the configured value

set-font-size[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-font-size "Link to this heading")
---------------------------------------------------------------------------------------------------

Fields are:

`size (required)`
The new font size in pts (a positive number). If absent is assumed to be zero which means reset to default.

`all (default: False)`
Boolean whether to change font size in the current window or all windows

`increment_op (optional)`
The string `+`, `-`, `*` or `/` to interpret size as an increment

set-spacing[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-spacing "Link to this heading")
-----------------------------------------------------------------------------------------------

Fields are:

`settings (required)`
An object mapping margins/paddings using canonical form {‘margin-top’: 50, ‘padding-left’: null} etc

`match_window (optional)`
Window to change paddings and margins in

`match_tab (default: None)`
Tab to change paddings and margins in

`all (default: False)`
Boolean indicating change paddings and margins everywhere or not

`configured (default: False)`
Boolean indicating whether to change the configured paddings and margins. Must be True if reset is True

set-tab-color[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-tab-color "Link to this heading")
---------------------------------------------------------------------------------------------------

Fields are:

`colors (required)`
An object mapping names to colors as 24-bit RGB integers. A color value of null indicates it should be unset.

`match (default: None)`
Which tab to change the color of

`self (default: False)`
Boolean indicating whether to use the tab of the window the command is run in

set-tab-title[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-tab-title "Link to this heading")
---------------------------------------------------------------------------------------------------

Fields are:

`title (required)`
The new title

`match (default: None)`
Which tab to change the title of

set-user-vars[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-user-vars "Link to this heading")
---------------------------------------------------------------------------------------------------

Fields are:

`var (optional)`
List of user variables of the form NAME=VALUE

`match (default: None)`
Which windows to change the title in

set-window-logo[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-window-logo "Link to this heading")
-------------------------------------------------------------------------------------------------------

Fields are:

`data (required)`
Chunk of PNG data, base64 encoded no more than 2048 bytes. Must send an empty chunk to indicate end of image. Or the special value `-` to indicate image must be removed.

`position (default: None)`
The logo position as a string, empty string means default

`alpha (default: -1.0)`
The logo alpha between `0` and `1`. `-1` means use default

`match (default: None)`
Which window to change the logo in

`self (default: False)`
Boolean indicating whether to act on the window the command is run in

set-window-title[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-window-title "Link to this heading")
---------------------------------------------------------------------------------------------------------

Fields are:

`title (optional)`
The new title

`match (default: None)`
Which windows to change the title in

`temporary (default: False)`
Boolean indicating if the change is temporary or permanent

signal-child[¶](https://sw.kovidgoyal.net/kitty/rc_protocol/#signal-child "Link to this heading")
-------------------------------------------------------------------------------------------------

Fields are:

`signals (required)`
The signals, a list of names, such as `SIGTERM`, `SIGKILL`, `SIGUSR1`, etc.

`match (default: None)`
Which windows to send the signals to

[Next The **launch** command](https://sw.kovidgoyal.net/kitty/launch/)[Previous Control kitty from scripts](https://sw.kovidgoyal.net/kitty/remote-control/)

 Copyright © 2025, Kovid Goyal 

 Made with [Furo](https://github.com/pradyunsg/furo)

[](https://github.com/kovidgoyal/kitty)

 On this page 

*   [The kitty remote control protocol](https://sw.kovidgoyal.net/kitty/rc_protocol/#)
    *   [Encrypted communication](https://sw.kovidgoyal.net/kitty/rc_protocol/#encrypted-communication)
    *   [Async and streaming requests](https://sw.kovidgoyal.net/kitty/rc_protocol/#async-and-streaming-requests)
    *   [action](https://sw.kovidgoyal.net/kitty/rc_protocol/#action)
    *   [close-tab](https://sw.kovidgoyal.net/kitty/rc_protocol/#close-tab)
    *   [close-window](https://sw.kovidgoyal.net/kitty/rc_protocol/#close-window)
    *   [create-marker](https://sw.kovidgoyal.net/kitty/rc_protocol/#create-marker)
    *   [detach-tab](https://sw.kovidgoyal.net/kitty/rc_protocol/#detach-tab)
    *   [detach-window](https://sw.kovidgoyal.net/kitty/rc_protocol/#detach-window)
    *   [disable-ligatures](https://sw.kovidgoyal.net/kitty/rc_protocol/#disable-ligatures)
    *   [env](https://sw.kovidgoyal.net/kitty/rc_protocol/#env)
    *   [focus-tab](https://sw.kovidgoyal.net/kitty/rc_protocol/#focus-tab)
    *   [focus-window](https://sw.kovidgoyal.net/kitty/rc_protocol/#focus-window)
    *   [get-colors](https://sw.kovidgoyal.net/kitty/rc_protocol/#get-colors)
    *   [get-text](https://sw.kovidgoyal.net/kitty/rc_protocol/#get-text)
    *   [goto-layout](https://sw.kovidgoyal.net/kitty/rc_protocol/#goto-layout)
    *   [kitten](https://sw.kovidgoyal.net/kitty/rc_protocol/#kitten)
    *   [last-used-layout](https://sw.kovidgoyal.net/kitty/rc_protocol/#last-used-layout)
    *   [launch](https://sw.kovidgoyal.net/kitty/rc_protocol/#launch)
    *   [load-config](https://sw.kovidgoyal.net/kitty/rc_protocol/#load-config)
    *   [ls](https://sw.kovidgoyal.net/kitty/rc_protocol/#ls)
    *   [new-window](https://sw.kovidgoyal.net/kitty/rc_protocol/#new-window)
    *   [remove-marker](https://sw.kovidgoyal.net/kitty/rc_protocol/#remove-marker)
    *   [resize-os-window](https://sw.kovidgoyal.net/kitty/rc_protocol/#resize-os-window)
    *   [resize-window](https://sw.kovidgoyal.net/kitty/rc_protocol/#resize-window)
    *   [run](https://sw.kovidgoyal.net/kitty/rc_protocol/#run)
    *   [scroll-window](https://sw.kovidgoyal.net/kitty/rc_protocol/#scroll-window)
    *   [select-window](https://sw.kovidgoyal.net/kitty/rc_protocol/#select-window)
    *   [send-key](https://sw.kovidgoyal.net/kitty/rc_protocol/#send-key)
    *   [send-text](https://sw.kovidgoyal.net/kitty/rc_protocol/#send-text)
    *   [set-background-image](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-background-image)
    *   [set-background-opacity](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-background-opacity)
    *   [set-colors](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-colors)
    *   [set-enabled-layouts](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-enabled-layouts)
    *   [set-font-size](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-font-size)
    *   [set-spacing](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-spacing)
    *   [set-tab-color](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-tab-color)
    *   [set-tab-title](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-tab-title)
    *   [set-user-vars](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-user-vars)
    *   [set-window-logo](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-window-logo)
    *   [set-window-title](https://sw.kovidgoyal.net/kitty/rc_protocol/#set-window-title)
    *   [signal-child](https://sw.kovidgoyal.net/kitty/rc_protocol/#signal-child)


