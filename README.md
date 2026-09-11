# <img src="icon.png" width="44" height="44" valign="middle" alt="OpenNox Icon" /> OpenNox for Android

Native Android port of **OpenNox**, the open-source community engine for Westwood Studios' classic ARPG **Nox** (2000).

> [!NOTE]
> This is an unofficial fan port. No copyrighted game assets are included. You must provide your own Nox data files (e.g. from GOG).

---

## 📥 Download & Support

* 📦 **Prebuilt APK:** [Download on Boosty](https://boosty.to/gignorie)
* ☕ **Boosty:** [boosty.to/gignorie](https://boosty.to/gignorie)
* 💎 **Crypto Donations (USDT):**
  * **USDT (TON):** `UQBPWSKGhTcl8rEXB2p7QR0xPHDZ5Y9Wn1T7YsnpvZP8yFMX`
  * **USDT (TRC20):** `TC5QTt58Z5d39YiffBy1J5d6CgQSw3pGuV`

---

## Overview

This repository contains the full Android build environment, native SDL2/OpenAL integration, touch controls, and packaging pipeline for running OpenNox natively on Android devices.

* **Target Architecture:** `armeabi-v7a` (32-bit ARM with VFPv3/NEON). Runs natively on ARMv7 and on 64-bit ARMv8/ARMv9 processors with 32-bit execution state (AArch32).
* **Minimum Android Version:** Android 7.0 (API level 24).
* **Target Android Version:** Android 14 (API level 34).
* **Graphics:** OpenGL ES 2.0 via SDL2.
* **Audio:** OpenAL Soft with OpenSL ES audio output backend.

---

## Directory Structure

```
opennox-android/
├── android-libs/          # Prebuilt native shared libraries (.so)
│   ├── armeabi-v7a/       # libSDL2.so, libopenal.so, libc++_shared.so, libmain.so
│   └── arm64-v8a/         # 64-bit dependencies
├── android-project/       # Android APK project source & packaging
│   ├── build_apk.sh       # Build script (aapt2, javac, d8, zipalign, apksigner)
│   ├── src/main/          # AndroidManifest.xml, SDLActivity, icons, strings
│   └── lib/               # Native libraries bundled into the APK
├── libs/                  # OpenNox Go libraries (input, seat, net, math)
│   └── client/seat/sdl/   # SDL seat & Android on-screen touch overlay
├── pkgconfig/             # pkg-config files for SDL2 and OpenAL (armv7a)
├── src/                   # OpenNox core engine source (Go + legacy C)
│   ├── CHANGELOG.md       # OpenNox & Android port changelog
│   └── src/cmd/opennox/   # Main entry point compiled into libmain.so
└── opennox-armeabi-v7a.apk# Final signed Android application package
```

---

## How to Build

### 1. Compile `libmain.so` (Go + CGO)

Run the following command inside `/root/opennox-android/src/src`:

```bash
cd /root/opennox-android/src/src && \
PKG_CONFIG_PATH=/root/opennox-android/pkgconfig \
CGO_CFLAGS_ALLOW=".*" \
CGO_LDFLAGS_ALLOW=".*" \
CGO_ENABLED=1 \
CC=/usr/local/bin/armv7a-linux-androideabi24-clang \
CGO_CFLAGS="-I/root/SDL2/include -I/root/SDL2/include/SDL2 -I/root/openal-soft/include -Wno-unknown-warning-option -fsigned-char" \
CGO_LDFLAGS="-L/root/opennox-android/android-libs/armeabi-v7a -lSDL2 -lopenal -llog -lGLESv2 -ldl" \
GOOS=android GOARCH=arm GOARM=7 \
go build -buildmode=c-shared -o /root/opennox-android/android-libs/armeabi-v7a/libmain.so -v ./cmd/opennox
```

### 2. Assemble and Sign the APK

Copy the newly compiled `libmain.so` into the APK project and execute `build_apk.sh`:

```bash
cp -f /root/opennox-android/android-libs/armeabi-v7a/libmain.so /root/opennox-android/android-project/lib/armeabi-v7a/libmain.so
/root/opennox-android/android-project/build_apk.sh
```

The signed APK will be generated at `/root/opennox-android/opennox-armeabi-v7a.apk` and copied automatically to `/sdcard/Download/opennox-armeabi-v7a.apk`.

---

## Game Installation & Data Setup

1. **Install the APK:**
   Download the latest prebuilt `opennox-armeabi-v7a.apk` from [Boosty](https://boosty.to/gignorie) (or compile it following the steps above) and install it on your Android device.

2. **Game Assets:**
   *Note: This is an unofficial fan port. No copyrighted game assets are included. You must provide your own Nox data files (e.g. from GOG).*

   Copy the original game data directory (containing maps, audio, modifier, video, etc.) to:
   ```
   /sdcard/opennox/
   ```
   *Required files include:* `NOX.EXE` or `game.exe`, `AUDIO/`, `VIDEO/`, `MAPS/`, `MODIFIER/`.

3. **Storage Permission:**
   On Android 11+ (API 30+), grant the **"All files access"** (`MANAGE_EXTERNAL_STORAGE`) permission when prompted, or enable it under *Settings → Apps → OpenNox → Permissions*.

---

## Controls & Touch Overlay

OpenNox Android features a built-in virtual touch overlay designed for touchscreen play:

* **Movement:** Virtual analog stick on the bottom-left area of the screen.
* **Actions:**
  * `[A]` — Primary attack / interact (Left Mouse Button).
  * `[J]` — Jump (Spacebar).
  * `[1]` to `[5]` — Quick-cast spells or abilities.
  * `[INV]` — Open / close inventory (`I`).
  * `[BOOK]` — Open / close spellbook (`B`).
  * `[ESC]` — Open menu / pause game (`Esc`).
* **Mouse Cursor / GUI:**
  * In menus and dialogs, the touch interface seamlessly transitions into direct touchscreen mouse pointing.
  * Tapping the screen moves the cursor and activates GUI controls.
* **Auto-Reset & Preemption:**
  * Multi-touch state machine automatically preempts stale touch IDs to prevent stuck clicks.
  * Touches are automatically reset upon menu transitions and window focus loss.

---

## Diagnostics & Logs

If an issue or crash occurs, diagnostic logs are written directly to the game directory on external storage:

* **Engine log:** `/sdcard/opennox/opennox.log`
* **Crash backtrace:** `/sdcard/opennox/crash.log`
