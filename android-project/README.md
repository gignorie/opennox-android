# OpenNox Android Packaging Project

This directory contains the standalone Android application packaging pipeline for OpenNox.

---

## Architecture

The project builds the final APK without relying on Android Studio or Gradle, using lightweight native Android build tools:
* `aapt2` — Compiles and links Android resources (manifest, drawables, strings).
* `javac` — Compiles the Java wrapper classes (`SDLActivity`, custom entry points).
* `d8` — Desugars and compiles Java `.class` files into Dalvik Executable (`classes.dex`).
* `zip` / `zipalign` — Packages DEX and `.so` native libraries aligned to 4-byte boundaries.
* `apksigner` — Signs the APK with APK Signature Scheme v2 & v3 using the debug keystore.

---

## Files

* `build_apk.sh` — The automated packaging and signing script.
* `src/main/AndroidManifest.xml` — Android application manifest (permissions, screen orientation, activities).
* `src/main/java/org/libsdl/app/SDLActivity.java` — SDL2 Android activity lifecycle handler.
* `src/main/res/` — Android resources (app icons, name, strings).
* `lib/armeabi-v7a/` — Bundled native libraries:
  * `libmain.so` (OpenNox engine)
  * `libSDL2.so` (SDL2 runtime)
  * `libopenal.so` (OpenAL Soft audio backend)
  * `libc++_shared.so` (LLVM C++ standard library)

---

## How to Package

```bash
cd /root/opennox-android/android-project
./build_apk.sh
```

Output APK will be written to:
* `/root/opennox-android/opennox-armeabi-v7a.apk`
* `/sdcard/Download/opennox-armeabi-v7a.apk`
