; 脚本由 Inno Setup 脚本向导 生成！
; 有关创建 Inno Setup 脚本文件的详细资料请查阅帮助文档！

#define MyAppName "YST_Unlock"
#define MyAppVersion "1.1"
#define MyAppPublisher "JohnRey"

[Setup]
; 注: AppId的值为单独标识该应用程序。
; 不要为其他安装程序使用相同的AppId值。
; (生成新的GUID，点击 工具|在IDE中生成GUID。)
AppId={{9B2DDDA2-0130-469E-9A9C-F9F33F358625}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
;AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={pf}\{#MyAppName}
DefaultGroupName={#MyAppName}
OutputBaseFilename=YST_Unlock
Compression=lzma
SolidCompression=yes

[Languages]
Name: "chinesesimp"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "E:\MyProject\YiSaiTongUnLock\UnlockAll\UnlockAll.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "E:\MyProject\YiSaiTongUnLock\UnlockAll\wps.exe"; DestDir: "{app}"; Flags: ignoreversion
; 注意: 不要在任何共享系统文件上使用“Flags: ignoreversion”

[Registry]
Root: HKCR; Subkey: "Directory\shell\Unlock"; ValueType: string; ValueName: ""; ValueData: "Unlock"; Flags: uninsdeletekey
Root: HKCR; Subkey: "Directory\shell\Unlock\command"; ValueType: string; ValueName: ""; ValueData: """{app}\UnlockAll.exe"" ""%1"""; Flags: uninsdeletekey
Root: HKCR; Subkey: "*\shell\Unlock"; ValueType: string; ValueName: ""; ValueData: "Unlock"; Flags: uninsdeletekey
Root: HKCR; Subkey: "*\shell\Unlock\command"; ValueType: string; ValueName: ""; ValueData: """{app}\UnlockAll.exe"" ""%1"""; Flags: uninsdeletekey
