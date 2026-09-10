#!/usr/bin/env bash
# OpenNox Android Launcher script

export DISPLAY="${DISPLAY:-:0}"
export GALLIUM_DRIVER="${GALLIUM_DRIVER:-virpipe}"
export MESA_GL_VERSION_OVERRIDE="${MESA_GL_VERSION_OVERRIDE:-3.3}"

# Ensure VirGL socket symlink exists inside PRoot if available in Termux
TERMUX_VIRGL_SOCK="/data/data/com.termux/files/usr/tmp/.virgl_test"
if [ -S "$TERMUX_VIRGL_SOCK" ] && [ ! -S "/tmp/.virgl_test" ]; then
    ln -sf "$TERMUX_VIRGL_SOCK" /tmp/.virgl_test
fi

GAME_BIN="/root/opennox-android/src/src/opennox"
GAME_DATA="/root/opennox-android/nox-data"

echo "=== Starting OpenNox ==="
echo "Display: $DISPLAY"
echo "Driver:  $GALLIUM_DRIVER"
echo "Data:    $GAME_DATA"

exec "$GAME_BIN" -data "$GAME_DATA" -window -nolimit "$@"
