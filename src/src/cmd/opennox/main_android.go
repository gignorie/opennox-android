//go:build android

package main

/*
#define SDL_MAIN_HANDLED
#include <stdlib.h>
#include <stdio.h>
#include <signal.h>
#include <ucontext.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>
#include <dlfcn.h>
#include <SDL2/SDL.h>
#include <android/log.h>

static void android_log(int prio, const char* tag, const char* msg) {
	__android_log_print(prio, tag, "%s", msg);
}

static void raw_crash_handler(int sig, siginfo_t* info, void* ucontext) {
	char buf[512];
	ucontext_t* uc = (ucontext_t*)ucontext;
	uintptr_t pc = 0, lr = 0, fault = (uintptr_t)info->si_addr;
#if defined(__arm__)
	pc = uc->uc_mcontext.arm_pc;
	lr = uc->uc_mcontext.arm_lr;
#elif defined(__aarch64__)
	pc = uc->uc_mcontext.pc;
	lr = uc->uc_mcontext.regs[30];
#endif
	const char* sym_pc = "?";
	const char* file_pc = "?";
	uintptr_t rel_pc = 0;
	Dl_info dli_pc;
	if (dladdr((void*)pc, &dli_pc)) {
		if (dli_pc.dli_sname) sym_pc = dli_pc.dli_sname;
		if (dli_pc.dli_fname) file_pc = dli_pc.dli_fname;
		if (dli_pc.dli_fbase) rel_pc = pc - (uintptr_t)dli_pc.dli_fbase;
	}

	const char* sym_lr = "?";
	const char* file_lr = "?";
	uintptr_t rel_lr = 0;
	Dl_info dli_lr;
	if (dladdr((void*)lr, &dli_lr)) {
		if (dli_lr.dli_sname) sym_lr = dli_lr.dli_sname;
		if (dli_lr.dli_fname) file_lr = dli_lr.dli_fname;
		if (dli_lr.dli_fbase) rel_lr = lr - (uintptr_t)dli_lr.dli_fbase;
	}

	int len = snprintf(buf, sizeof(buf),
		"\n\n*** FATAL NATIVE CRASH: signal %d (si_code=%d), fault addr=%p, PC=%p (+0x%lx in %s, sym=%s), LR=%p (+0x%lx in %s, sym=%s) ***\n\n",
		sig, info->si_code, (void*)fault, (void*)pc, (unsigned long)rel_pc, file_pc, sym_pc, (void*)lr, (unsigned long)rel_lr, file_lr, sym_lr);

	__android_log_print(ANDROID_LOG_FATAL, "OpenNoxCrash", "%s", buf);

	int fd = open("/sdcard/opennox/crash.log", O_WRONLY | O_CREAT | O_APPEND, 0666);
	if (fd >= 0) {
		write(fd, buf, len);
		close(fd);
	}
	int fd2 = open("/sdcard/opennox/opennox.log", O_WRONLY | O_CREAT | O_APPEND, 0666);
	if (fd2 >= 0) {
		write(fd2, buf, len);
		close(fd2);
	}
	_exit(128 + sig);
}

static void install_crash_handlers(void) {
	struct sigaction sa;
	memset(&sa, 0, sizeof(sa));
	sa.sa_sigaction = raw_crash_handler;
	sa.sa_flags = SA_SIGINFO | SA_ONSTACK;
	sigaction(SIGSEGV, &sa, NULL);
	sigaction(SIGBUS, &sa, NULL);
	sigaction(SIGILL, &sa, NULL);
	sigaction(SIGFPE, &sa, NULL);
	sigaction(SIGABRT, &sa, NULL);
}
*/
import "C"
import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"syscall"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/opennox/opennox/v1"
)

func main() {}

var logFile *os.File

type syncFileWriter struct {
	f *os.File
}

func (s syncFileWriter) Write(p []byte) (int, error) {
	n, err := s.f.Write(p)
	_ = s.f.Sync()
	return n, err
}

func logToAndroidAndFile(prio C.int, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	cTag := C.CString("OpenNox")
	cMsg := C.CString(msg)
	C.android_log(prio, cTag, cMsg)
	C.free(unsafe.Pointer(cTag))
	C.free(unsafe.Pointer(cMsg))

	if logFile != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(logFile, "[%s] %s\n", timestamp, msg)
		logFile.Sync()
	}
}

func checkNoxDir(dir string) bool {
	if dir == "" {
		return false
	}
	if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
		checkFiles := []string{
			"gamedata.bin", "GAMEDATA.BIN",
			"thing.bin", "THING.BIN",
			"thing.bag", "THING.BAG", "Thing.bag",
			"nox.gxm", "NOX.GXM",
			"game.exe", "GAME.EXE",
		}
		for _, f := range checkFiles {
			if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
				return true
			}
		}
	}
	return false
}

//export SDL_main
func SDL_main(argc C.int, argv **C.char) C.int {
	C.SDL_SetMainReady()
	C.install_crash_handlers()

	// Try opening log file in accessible locations
	for _, logCandidate := range []string{
		"/sdcard/opennox/opennox.log",
		"/sdcard/opennox_debug.log",
		filepath.Join(sdl.AndroidGetExternalStoragePath(), "opennox.log"),
		filepath.Join(sdl.AndroidGetInternalStoragePath(), "opennox.log"),
	} {
		if dir := filepath.Dir(logCandidate); dir != "" {
			_ = os.MkdirAll(dir, 0755)
		}
		if f, err := os.OpenFile(logCandidate, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			logFile = f
			fd := int(f.Fd())
			_ = syscall.Dup2(fd, 1)
			_ = syscall.Dup2(fd, 2)
			os.Stdout = f
			os.Stderr = f
			slog.SetDefault(slog.New(slog.NewTextHandler(syncFileWriter{f: f}, &slog.HandlerOptions{Level: slog.LevelDebug})))
			break
		}
	}

	logToAndroidAndFile(C.ANDROID_LOG_INFO, "=== OpenNox SDL_main invoked from Android SDLActivity ===")

	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			errText := fmt.Sprintf("FATAL PANIC:\n%v\n\nStack:\n%s", r, stack)
			logToAndroidAndFile(C.ANDROID_LOG_ERROR, "%s", errText)
			sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "OpenNox Crash", errText, nil)
		}
		if logFile != nil {
			_ = logFile.Close()
		}
	}()

	var args []string
	if argv != nil && argc > 0 {
		cArgs := (*[1 << 20]*C.char)(unsafe.Pointer(argv))[:int(argc):int(argc)]
		for _, arg := range cArgs {
			if arg != nil {
				args = append(args, C.GoString(arg))
			}
		}
	}
	if len(args) == 0 {
		args = []string{"opennox"}
	}

	hasData := false
	for _, a := range args {
		if a == "-data" {
			hasData = true
			break
		}
	}

	if !hasData {
		var candidates []string
		if ext := sdl.AndroidGetExternalStoragePath(); ext != "" {
			candidates = append(candidates, ext, filepath.Join(ext, "nox"), filepath.Join(ext, "opennox"))
		}
		if intern := sdl.AndroidGetInternalStoragePath(); intern != "" {
			candidates = append(candidates, intern, filepath.Join(intern, "nox"), filepath.Join(intern, "opennox"))
		}
		candidates = append(candidates,
			"/sdcard/OpenNox",
			"/sdcard/opennox",
			"/sdcard/nox",
			"/storage/emulated/0/OpenNox",
			"/storage/emulated/0/opennox",
			"/storage/emulated/0/nox",
			"/sdcard/Android/data/org.libsdl.app/files",
			"/sdcard/Android/data/com.opennox/files",
			"/sdcard/Android/data/org.libsdl.app/files/nox",
			"/sdcard/Android/data/com.opennox/files/nox",
		)

		var chosen string
		logToAndroidAndFile(C.ANDROID_LOG_INFO, "Searching for Nox data directory...")
		for attempt := 0; attempt < 20; attempt++ {
			for _, dir := range candidates {
				if checkNoxDir(dir) {
					chosen = dir
					break
				}
			}
			if chosen != "" {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		if chosen != "" {
			logToAndroidAndFile(C.ANDROID_LOG_INFO, "Found Nox data directory: %s", chosen)
			_ = os.Chdir(chosen)
			args = append(args, "-data", chosen)
		} else {
			errMsg := fmt.Sprintf("Nox game data files (gamedata.bin, thing.bin) not found!\n\nPlease make sure game data is in /sdcard/opennox\nand 'All files access' permission is granted.")
			logToAndroidAndFile(C.ANDROID_LOG_ERROR, "%s", errMsg)
			sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "OpenNox Data Not Found", errMsg, nil)
			return 1
		}
	}

	logToAndroidAndFile(C.ANDROID_LOG_INFO, "Running opennox with args: %v", args)
	slog.Info("Running opennox with args", "args", args)
	if err := opennox.RunArgs(args); err != nil {
		errMsg := fmt.Sprintf("opennox.RunArgs terminated with error:\n%v", err)
		logToAndroidAndFile(C.ANDROID_LOG_ERROR, "%s", errMsg)
		sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_ERROR, "OpenNox Error", errMsg, nil)
		os.Exit(1)
		return 1
	}
	logToAndroidAndFile(C.ANDROID_LOG_INFO, "opennox.RunArgs finished normally")
	os.Exit(0)
	return 0
}
