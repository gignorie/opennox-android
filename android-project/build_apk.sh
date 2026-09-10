#!/bin/bash
set -e

PROJECT_DIR="/root/opennox-android/android-project"
SDK_JAR="/root/android-sdk/platforms/android-34/android.jar"
BUILD_TOOLS="/root/android-sdk/build-tools/34.0.0"
KEYSTORE="/root/.android/debug.keystore"
OUTPUT_APK="/root/opennox-android/opennox-armeabi-v7a.apk"

cd "$PROJECT_DIR"

echo "=== 1. Clean build directories ==="
rm -rf build/compiled_res build/compiled_res.zip build/gen build/classes build/dex build/*.apk
mkdir -p build/gen build/classes build/dex

echo "=== 2. Compile Android Resources (aapt2) ==="
/usr/bin/aapt2 compile --dir src/main/res -o build/compiled_res.zip

echo "=== 3. Link Android Resources and create R.java (aapt2) ==="
/usr/bin/aapt2 link -I "$SDK_JAR" \
    --manifest src/main/AndroidManifest.xml \
    --java build/gen \
    -o build/unaligned_base.apk \
    build/compiled_res.zip

echo "=== 4. Compile Java sources (javac) ==="
javac -encoding UTF-8 -cp "$SDK_JAR" \
    -sourcepath src/main/java:build/gen \
    -d build/classes \
    $(find src/main/java build/gen -name "*.java")

echo "=== 5. Convert classes to DEX (d8) ==="
"$BUILD_TOOLS/d8" --lib "$SDK_JAR" \
    --output build/dex \
    $(find build/classes -name "*.class")

echo "=== 6. Package DEX and native libraries ==="
cd build/dex
zip -0 -u "$PROJECT_DIR/build/unaligned_base.apk" classes.dex
cd "$PROJECT_DIR"
zip -0 -u -r build/unaligned_base.apk lib/

echo "=== 7. Zipalign APK ==="
zipalign -p -f 4 build/unaligned_base.apk build/aligned.apk

echo "=== 8. Sign APK with apksigner ==="
apksigner sign --ks "$KEYSTORE" \
    --ks-pass pass:android \
    --key-pass pass:android \
    --out "$OUTPUT_APK" \
    build/aligned.apk

echo "=== 9. Verify signed APK ==="
apksigner verify -v "$OUTPUT_APK"

echo "=== 10. Copy APK to SDCard locations ==="
cp -f "$OUTPUT_APK" /sdcard/Download/opennox-armeabi-v7a.apk
cp -f "$OUTPUT_APK" /sdcard/opennox-armeabi-v7a.apk

echo "SUCCESS! APK built and signed: $OUTPUT_APK"
