package com.opennox;

import android.Manifest;
import android.app.Activity;
import android.app.AlertDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.pm.PackageManager;
import android.graphics.Color;
import android.graphics.Typeface;
import android.net.Uri;
import android.os.Build;
import android.os.Bundle;
import android.os.Environment;
import android.os.Handler;
import android.os.Looper;
import android.provider.DocumentsContract;
import android.provider.Settings;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.ProgressBar;
import android.widget.TextView;
import android.widget.Toast;

import java.io.BufferedInputStream;
import java.io.File;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

public class LauncherActivity extends Activity {

    private static final int REQUEST_STORAGE_PERMISSIONS = 101;
    private static final int REQUEST_MANAGE_STORAGE = 102;
    private static final int REQUEST_PICK_FOLDER = 103;
    private static final int REQUEST_PICK_ZIP = 104;

    private static final String PREFS_NAME = "opennox_launcher";
    private static final String PREF_CUSTOM_DIR = "custom_game_dir";

    private TextView mTvTitle;
    private TextView mTvSubtitle;
    private TextView mTvStatusHeader;
    private TextView mTvStatusMessage;
    private TextView mTvSelectedPath;

    private LinearLayout mLayoutProgress;
    private TextView mTvProgressInfo;
    private ProgressBar mProgressBar;
    private TextView mTvProgressPercent;

    private Button mBtnSelectFolder;
    private Button mBtnUnpackArchive;
    private Button mBtnCommunity;
    private Button mBtnLaunchGame;

    private File mCurrentGameDir = null;
    private boolean mIsUnpacking = false;
    private final Handler mHandler = new Handler(Looper.getMainLooper());

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_launcher);

        initViews();
        applyFonts();
        setupListeners();

        // Check if game files exist and we are not in explicit setup mode
        boolean forceSetup = getIntent().getBooleanExtra("setup_mode", false);
        if (checkStoragePermissionsGranted()) {
            File dir = findGameDirectory();
            if (dir != null && !forceSetup) {
                launchGame(dir);
                return;
            }
        }

        refreshGameStatus();
        requestStoragePermissionsIfNeeded();
    }

    @Override
    protected void onResume() {
        super.onResume();
        if (!mIsUnpacking) {
            refreshGameStatus();
        }
    }

    private void initViews() {
        mTvTitle = findViewById(R.id.tv_title);
        mTvSubtitle = findViewById(R.id.tv_subtitle);
        mTvStatusHeader = findViewById(R.id.tv_status_header);
        mTvStatusMessage = findViewById(R.id.tv_status_message);
        mTvSelectedPath = findViewById(R.id.tv_selected_path);

        mLayoutProgress = findViewById(R.id.layout_progress);
        mTvProgressInfo = findViewById(R.id.tv_progress_info);
        mProgressBar = findViewById(R.id.progress_bar);
        mTvProgressPercent = findViewById(R.id.tv_progress_percent);

        mBtnSelectFolder = findViewById(R.id.btn_select_folder);
        mBtnUnpackArchive = findViewById(R.id.btn_unpack_archive);
        mBtnCommunity = findViewById(R.id.btn_community);
        mBtnLaunchGame = findViewById(R.id.btn_launch_game);
    }

    private void applyFonts() {
        try {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                Typeface tf = getResources().getFont(R.font.morris_roman);
                if (tf != null) {
                    if (mTvTitle != null) mTvTitle.setTypeface(tf);
                    if (mTvSubtitle != null) mTvSubtitle.setTypeface(tf);
                    if (mTvStatusHeader != null) mTvStatusHeader.setTypeface(tf, Typeface.BOLD);
                    if (mBtnSelectFolder != null) mBtnSelectFolder.setTypeface(tf, Typeface.BOLD);
                    if (mBtnUnpackArchive != null) mBtnUnpackArchive.setTypeface(tf, Typeface.BOLD);
                    if (mBtnCommunity != null) mBtnCommunity.setTypeface(tf, Typeface.BOLD);
                    if (mBtnLaunchGame != null) mBtnLaunchGame.setTypeface(tf, Typeface.BOLD);
                }
            }
        } catch (Exception ignored) {}
    }

    private void setupListeners() {
        mBtnSelectFolder.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                if (ensureStoragePermissions()) {
                    openFolderPicker();
                }
            }
        });

        mBtnUnpackArchive.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                if (ensureStoragePermissions()) {
                    openZipPicker();
                }
            }
        });

        mBtnCommunity.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                showCommunityAndGuidesDialog();
            }
        });

        mBtnLaunchGame.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                if (mCurrentGameDir != null && checkGameFiles(mCurrentGameDir)) {
                    launchGame(mCurrentGameDir);
                } else {
                    File dir = findGameDirectory();
                    if (dir != null) {
                        launchGame(dir);
                    } else {
                        Toast.makeText(LauncherActivity.this, "Файлы игры не найдены! Укажите папку или распакуйте архив.", Toast.LENGTH_LONG).show();
                    }
                }
            }
        });
    }

    private boolean checkStoragePermissionsGranted() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
            return Environment.isExternalStorageManager();
        } else if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
            return checkSelfPermission(Manifest.permission.READ_EXTERNAL_STORAGE) == PackageManager.PERMISSION_GRANTED &&
                   checkSelfPermission(Manifest.permission.WRITE_EXTERNAL_STORAGE) == PackageManager.PERMISSION_GRANTED;
        }
        return true;
    }

    private boolean ensureStoragePermissions() {
        if (!checkStoragePermissionsGranted()) {
            requestStoragePermissionsIfNeeded();
            return false;
        }
        return true;
    }

    private void requestStoragePermissionsIfNeeded() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
            if (!Environment.isExternalStorageManager()) {
                new AlertDialog.Builder(this)
                    .setTitle("Разрешение на доступ к файлам")
                    .setMessage("Для работы OpenNox и загрузки игровых архивов требуется предоставить 'Доступ ко всем файлам'.")
                    .setPositiveButton("Настройки", new DialogInterface.OnClickListener() {
                        @Override
                        public void onClick(DialogInterface dialog, int which) {
                            try {
                                Intent intent = new Intent(Settings.ACTION_MANAGE_APP_ALL_FILES_ACCESS_PERMISSION);
                                intent.setData(Uri.parse("package:" + getPackageName()));
                                startActivityForResult(intent, REQUEST_MANAGE_STORAGE);
                            } catch (Exception e) {
                                try {
                                    Intent intent = new Intent(Settings.ACTION_MANAGE_ALL_FILES_ACCESS_PERMISSION);
                                    startActivityForResult(intent, REQUEST_MANAGE_STORAGE);
                                } catch (Exception ignored) {}
                            }
                        }
                    })
                    .setNegativeButton("Позже", null)
                    .show();
            }
        } else if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
            if (checkSelfPermission(Manifest.permission.READ_EXTERNAL_STORAGE) != PackageManager.PERMISSION_GRANTED ||
                checkSelfPermission(Manifest.permission.WRITE_EXTERNAL_STORAGE) != PackageManager.PERMISSION_GRANTED) {
                requestPermissions(new String[]{
                    Manifest.permission.READ_EXTERNAL_STORAGE,
                    Manifest.permission.WRITE_EXTERNAL_STORAGE
                }, REQUEST_STORAGE_PERMISSIONS);
            }
        }
    }

    private void refreshGameStatus() {
        mCurrentGameDir = findGameDirectory();
        if (mCurrentGameDir != null) {
            mTvStatusMessage.setText("Файлы игры обнаружены!\nУстановка готова к запуску.");
            mTvStatusMessage.setTextColor(Color.parseColor("#7ecb5c"));
            mTvSelectedPath.setText(mCurrentGameDir.getAbsolutePath());
            mTvSelectedPath.setTextColor(Color.parseColor("#caa048"));
            mBtnLaunchGame.setEnabled(true);
        } else {
            mTvStatusMessage.setText("Файлы игры не найдены.\n\nДля запуска требуются ресурсы Westwood Studios Nox (thing.bag, video.bag, Audio.bag, game.exe).\nУкажите папку с установленной игрой или выберите zip-архив.");
            mTvStatusMessage.setTextColor(Color.parseColor("#e0c7a8"));
            File defaultDir = new File(Environment.getExternalStorageDirectory(), "OpenNox");
            mTvSelectedPath.setText("Ожидается: " + defaultDir.getAbsolutePath() + "/");
            mTvSelectedPath.setTextColor(Color.parseColor("#8f7f68"));
            mBtnLaunchGame.setEnabled(false);
        }
    }

    private File findGameDirectory() {
        SharedPreferences prefs = getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE);
        String customPath = prefs.getString(PREF_CUSTOM_DIR, null);
        if (customPath != null) {
            File customDir = new File(customPath);
            if (checkGameFiles(customDir)) {
                return customDir;
            }
        }

        File ext = Environment.getExternalStorageDirectory();
        File[] candidates = new File[] {
            new File(ext, "OpenNox"),
            new File(ext, "opennox"),
            new File("/sdcard/OpenNox"),
            new File("/sdcard/opennox"),
            new File("/storage/emulated/0/OpenNox"),
            new File("/storage/emulated/0/opennox"),
            new File(ext, "Nox"),
            new File(ext, "nox"),
            getExternalFilesDir(null),
            new File(getExternalFilesDir(null), "nox")
        };

        for (File c : candidates) {
            if (c != null && checkGameFiles(c)) {
                return c;
            }
        }
        return null;
    }

    private boolean checkGameFiles(File dir) {
        if (dir == null || !dir.exists() || !dir.isDirectory()) {
            return false;
        }
        String[] essentialFiles = new String[] {
            "thing.bag", "THING.BAG", "Thing.bag",
            "gamedata.bin", "GAMEDATA.BIN",
            "thing.bin", "THING.BIN",
            "nox.gxm", "NOX.GXM",
            "game.exe", "GAME.EXE",
            "nox.exe", "NOX.EXE"
        };
        for (String name : essentialFiles) {
            File f = new File(dir, name);
            if (f.exists() && f.isFile()) {
                return true;
            }
        }
        return false;
    }

    private void openFolderPicker() {
        try {
            Intent intent = new Intent(Intent.ACTION_OPEN_DOCUMENT_TREE);
            intent.addFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION | Intent.FLAG_GRANT_WRITE_URI_PERMISSION);
            startActivityForResult(intent, REQUEST_PICK_FOLDER);
        } catch (Exception e) {
            Toast.makeText(this, "Не удалось открыть выбор папки: " + e.getMessage(), Toast.LENGTH_SHORT).show();
        }
    }

    private void openZipPicker() {
        try {
            Intent intent = new Intent(Intent.ACTION_OPEN_DOCUMENT);
            intent.addCategory(Intent.CATEGORY_OPENABLE);
            intent.setType("*/*");
            String[] mimeTypes = new String[]{
                "application/zip",
                "application/x-zip-compressed",
                "application/octet-stream",
                "*/*"
            };
            intent.putExtra(Intent.EXTRA_MIME_TYPES, mimeTypes);
            startActivityForResult(intent, REQUEST_PICK_ZIP);
        } catch (Exception e) {
            Toast.makeText(this, "Не удалось открыть выбор архива: " + e.getMessage(), Toast.LENGTH_SHORT).show();
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);

        if (requestCode == REQUEST_MANAGE_STORAGE) {
            refreshGameStatus();
            return;
        }

        if (resultCode != RESULT_OK || data == null) {
            return;
        }

        if (requestCode == REQUEST_PICK_FOLDER) {
            Uri treeUri = data.getData();
            if (treeUri != null) {
                try {
                    getContentResolver().takePersistableUriPermission(
                        treeUri,
                        Intent.FLAG_GRANT_READ_URI_PERMISSION | Intent.FLAG_GRANT_WRITE_URI_PERMISSION
                    );
                } catch (Exception ignored) {}

                String resolvedPath = resolvePathFromUri(treeUri);
                if (resolvedPath != null) {
                    File chosenDir = new File(resolvedPath);
                    if (checkGameFiles(chosenDir)) {
                        saveCustomGameDir(chosenDir.getAbsolutePath());
                        Toast.makeText(this, "Папка с игрой найдена и сохранена!", Toast.LENGTH_SHORT).show();
                    } else {
                        saveCustomGameDir(chosenDir.getAbsolutePath());
                        Toast.makeText(this, "Выбрана папка: " + resolvedPath + ", но файлы игры не найдены.", Toast.LENGTH_LONG).show();
                    }
                } else {
                    Toast.makeText(this, "Выбран URI: " + treeUri.getPath(), Toast.LENGTH_SHORT).show();
                }
                refreshGameStatus();
            }
        } else if (requestCode == REQUEST_PICK_ZIP) {
            Uri zipUri = data.getData();
            if (zipUri != null) {
                startUnpackZip(zipUri);
            }
        }
    }

    private String resolvePathFromUri(Uri treeUri) {
        try {
            String docId = DocumentsContract.getTreeDocumentId(treeUri);
            if (docId != null) {
                String[] parts = docId.split(":");
                String type = parts[0];
                String relativePath = parts.length > 1 ? parts[1] : "";
                if ("primary".equalsIgnoreCase(type)) {
                    return Environment.getExternalStorageDirectory().getAbsolutePath() + "/" + relativePath;
                } else {
                    return "/storage/" + type + "/" + relativePath;
                }
            }
        } catch (Exception ignored) {}
        return null;
    }

    private void saveCustomGameDir(String path) {
        SharedPreferences.Editor editor = getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE).edit();
        editor.putString(PREF_CUSTOM_DIR, path);
        editor.apply();
    }

    private void startUnpackZip(final Uri zipUri) {
        if (mIsUnpacking) {
            return;
        }
        mIsUnpacking = true;
        setButtonsEnabled(false);
        mLayoutProgress.setVisibility(View.VISIBLE);
        mProgressBar.setIndeterminate(true);
        mTvProgressInfo.setText("Подготовка к распаковке...");
        mTvProgressPercent.setText("");

        final File destDir = new File(Environment.getExternalStorageDirectory(), "OpenNox");
        if (!destDir.exists()) {
            destDir.mkdirs();
        }

        new Thread(new Runnable() {
            @Override
            public void run() {
                boolean success = false;
                String errorMsg = null;
                int extractedCount = 0;

                try {
                    InputStream is = getContentResolver().openInputStream(zipUri);
                    if (is == null) {
                        throw new Exception("Не удалось открыть файл архива");
                    }
                    BufferedInputStream bis = new BufferedInputStream(is, 65536);
                    ZipInputStream zis = new ZipInputStream(bis);

                    byte[] buffer = new byte[65536];
                    ZipEntry entry;
                    long lastUpdateTime = 0;

                    while ((entry = zis.getNextEntry()) != null) {
                        String name = entry.getName();
                        // Security check: Zip Slip vulnerability protection
                        File outFile = new File(destDir, name);
                        String canonicalDest = destDir.getCanonicalPath();
                        String canonicalOut = outFile.getCanonicalPath();
                        if (!canonicalOut.startsWith(canonicalDest)) {
                            zis.closeEntry();
                            continue;
                        }

                        if (entry.isDirectory()) {
                            outFile.mkdirs();
                        } else {
                            File parent = outFile.getParentFile();
                            if (parent != null && !parent.exists()) {
                                parent.mkdirs();
                            }
                            FileOutputStream fos = new FileOutputStream(outFile);
                            int len;
                            while ((len = zis.read(buffer)) > 0) {
                                fos.write(buffer, 0, len);
                            }
                            fos.flush();
                            fos.close();
                            extractedCount++;
                        }
                        zis.closeEntry();

                        long now = System.currentTimeMillis();
                        if (now - lastUpdateTime > 150) {
                            lastUpdateTime = now;
                            final String currentName = name;
                            final int count = extractedCount;
                            mHandler.post(new Runnable() {
                                @Override
                                public void run() {
                                    mProgressBar.setIndeterminate(false);
                                    mTvProgressInfo.setText("Извлечение: " + currentName);
                                    mTvProgressPercent.setText("Файлов: " + count);
                                }
                            });
                        }
                    }
                    zis.close();
                    success = true;

                    // Also create/sync opennox lowercase folder if needed
                    File opennoxLower = new File(Environment.getExternalStorageDirectory(), "opennox");
                    if (!opennoxLower.exists()) {
                        opennoxLower.mkdirs();
                    }
                } catch (Exception e) {
                    errorMsg = e.getMessage();
                }

                final boolean finalSuccess = success;
                final String finalError = errorMsg;
                final int finalCount = extractedCount;

                mHandler.post(new Runnable() {
                    @Override
                    public void run() {
                        mIsUnpacking = false;
                        mLayoutProgress.setVisibility(View.GONE);
                        setButtonsEnabled(true);

                        if (finalSuccess) {
                            saveCustomGameDir(destDir.getAbsolutePath());
                            Toast.makeText(LauncherActivity.this, "Распаковано файлов: " + finalCount + "! Игра готова к запуску.", Toast.LENGTH_LONG).show();
                            refreshGameStatus();
                        } else {
                            Toast.makeText(LauncherActivity.this, "Ошибка распаковки: " + (finalError != null ? finalError : "неизвестная ошибка"), Toast.LENGTH_LONG).show();
                            refreshGameStatus();
                        }
                    }
                });
            }
        }).start();
    }

    private void setButtonsEnabled(boolean enabled) {
        mBtnSelectFolder.setEnabled(enabled);
        mBtnUnpackArchive.setEnabled(enabled);
        mBtnCommunity.setEnabled(enabled);
        mBtnLaunchGame.setEnabled(enabled && (mCurrentGameDir != null));
    }

    private void showCommunityAndGuidesDialog() {
        AlertDialog.Builder builder = new AlertDialog.Builder(this);
        builder.setTitle("Сообщество OpenNox и гайды");
        builder.setMessage(
            "Добро пожаловать в порт OpenNox для Android!\n\n" +
            "• Для игры требуются файлы полной версии Nox (1.2 или GOG):\n" +
            "  - thing.bag, video.bag, Audio.bag\n" +
            "  - gamedata.bin, thing.bin, game.exe\n" +
            "  - папки Dialog, maps, window\n\n" +
            "• Поместите файлы в /sdcard/OpenNox/ или воспользуйтесь кнопкой 'Распаковать архив'.\n\n" +
            "Выберите ссылку для перехода:"
        );
        builder.setPositiveButton("Telegram канал", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                openUrl("https://t.me/opennox");
            }
        });
        builder.setNeutralButton("GitHub OpenNox", new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                openUrl("https://github.com/opennox/opennox");
            }
        });
        builder.setNegativeButton("Закрыть", null);
        builder.show();
    }

    private void openUrl(String url) {
        try {
            Intent intent = new Intent(Intent.ACTION_VIEW, Uri.parse(url));
            startActivity(intent);
        } catch (Exception e) {
            Toast.makeText(this, "Не удалось открыть браузер: " + e.getMessage(), Toast.LENGTH_SHORT).show();
        }
    }

    private void launchGame(File dir) {
        Intent intent = new Intent(this, MainActivity.class);
        if (dir != null) {
            intent.putExtra("nox_data_dir", dir.getAbsolutePath());
        }
        startActivity(intent);
        finish();
    }
}
