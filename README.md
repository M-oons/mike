# mike

Mike is a lightweight Windows utility for global audio output muting/unmuting. It supports customizable hotkeys, system tray integration, and works with both standard Windows audio devices and the Voicemeeter virtual audio mixer.

## Install / Build
```batch
> make install
> make
```

## Configuration

The configuration file for Mike can be found in your user AppData roaming folder: `%APPDATA%/mike/config.json`. When Mike is launched for the first time, the configuration file is created with its default values. Each configuration section is outlined below.

| Key          | Description                    | Type   |
|--------------|--------------------------------|--------|
| `hotkeys`    | List of hotkey definitions     | Array  |
| `sounds`     | Sound output settings          | Object |
| `controller` | Audio controller configuration | Object |

---

### `hotkeys`

Defines one or more hotkey bindings.

| Key     | Description                          | Type    | Possible values                                     |
|---------|--------------------------------------|---------|-----------------------------------------------------|
| action  | The action performed by the hotkey   | String  | `mute`, `unmute`, `toggle`                          |
| key     | The main key for the hotkey          | String  | [See available hotkey keys](#available-hotkey-keys) |
| ctrl    | Whether `ctrl` must be held          | Boolean | `true`, `false`                                     |
| shift   | Whether `shift` must be held         | Boolean | `true`, `false`                                     |
| alt     | Whether `alt` must be held           | Boolean | `true`, `false`                                     |
| win     | Whether the Windows key must be held | Boolean | `true`, `false`                                     |

#### Available hotkey keys:

`a` `b` `c` `d` `e` `f` `g` `h` `i` `j` `k` `l` `m` `n` `o` `p` `q` `r` `s` `t` `u` `v` `w` `x` `y` `z`<br>
`0` `1` `2` `3` `4` `5` `6` `7` `8` `9`<br>
`` ` `` `'` `-` `=` `#` `,` `.` `;` `/` `\` `[` `]`<br>
`numpad0` `numpad1` `numpad2` `numpad3` `numpad4` `numpad5` `numpad6` `numpad7` `numpad8` `numpad9`<br>
`numpad+` `numpad-` `numpad*` `numpad/` `numpad.`<br>
`f1` `f2` `f3` `f4` `f5` `f6` `f7` `f8` `f9` `f10` `f11` `f12`<br>
`f13` `f14` `f15` `f16` `f17` `f18` `f19` `f20` `f21` `f22` `f23` `f24`<br>

---

### `sounds`

Controls sound output for actions.

| Key     | Description                     | Type    | Possible values | Default |
|---------|---------------------------------|---------|-----------------|---------|
| enabled | Whether sound output is enabled | Boolean | `true`, `false` | `true`  |
| volume  | Playback volume (%)             | Number  | -               | `100`   |

---

### `controller`

Configures which audio controller to use and specific options.

| Key         | Description                                            | Type    | Possible values          | Default   |
|-------------|--------------------------------------------------------|---------|--------------------------|-----------|
| type        | Which audio backend to use                             | String  | `windows`, `voicemeeter` | `windows` |
| windows     | Windows specific controller options (currently unused) | Object  | -                        | -         |
| voicemeeter | Voicemeeter specific controller options                | Object  | -                        | -         |

### `controller.windows`

Windows specific controller options (currently unused).

### `controller.voicemeeter`

Voicemeeter specific controller options.

| Key           | Description                                 | Type   | Possible values                                    | Default                                                         |
|---------------|---------------------------------------------|--------|----------------------------------------------------|-----------------------------------------------------------------|
| remoteDLLPath | Full path to the Voicemeeter Remote API DLL | String | -                                                  | `C:/Program Files (x86)/VB/Voicemeeter/VoicemeeterRemote64.dll` |
| output        | Voicemeeter virtual output number           | Number | `1` (Virtual output B1)<br>`2` (Virtual output B2) | `1`                                                             |

---