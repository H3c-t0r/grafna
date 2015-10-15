#define MyAppID "{5FFA65A5-D4CF-4E26-9AC0-1615E3895B1E}"
#define MyAppName "AdRem GrafCrunch Server"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "AdRem Software, Inc. New York, NY"
#define MyAppURL "http://www.adremsoft.com/"
#define MyAppIcon "icon.ico"
#define MyAppGroupName "AdRem GrafCrunch"
#define MyAppUninstallKey "SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\" + MyAppID + "_is1"

#define LICENSE "..\..\LICENSE.md"
#define NOTICE "..\..\NOTICE.md"

#define GrafCrunchProgramData "{%programdata}\AdRem\GrafCrunch"

#define ConfigINIFile "\conf\custom.ini"
#define ConfigINI "{app}" + ConfigINIFile
#define PathConfigSection "paths"

#define GrafCrunchServerSection "server"
#define DefaultGrafCrunchServerDomain "localhost"
#define DefaultGrafCrunchServerPort "3000"

#define NetCrunchServerConfigSection "netcrunch-server"
#define DefaultNetCrunchServerAddress "localhost"
#define DefaultNetCrunchServerPort "12009"
#define DefaultNetCrunchServerUser "GrafCrunch"

#define NetCrunchServerKey "SOFTWARE\AdRem\NetCrunch\9.0\NCServer\Options\ServerConfiguration"

[Setup]
AppId={{#MyAppID}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppName}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
LicenseFile={#LICENSE}
DefaultDirName={pf64}\AdRem\GrafCrunch
DefaultGroupName={#MyAppGroupName}
OutputDir=release
OutputBaseFilename=GCServerSetup
SetupIconFile={#MyAppIcon}
UninstallDisplayIcon={#MyAppIcon}
Compression=lzma
SolidCompression=yes
ArchitecturesAllowed=x64
ArchitecturesInstallIn64BitMode=x64

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Registry]
Root: HKLM64; Subkey: "Software\AdRem"; Flags: uninsdeletekeyifempty
Root: HKLM64; Subkey: "Software\AdRem\GrafCrunch"; Flags: uninsdeletekey
Root: HKLM64; Subkey: "Software\AdRem\GrafCrunch\1.0"; Flags: uninsdeletekey
Root: HKLM64; Subkey: "Software\AdRem\GrafCrunch\1.0"; ValueType: string; ValueName: "AppFolder"; ValueData: "{app}"
Root: HKLM64; Subkey: "Software\AdRem\GrafCrunch\1.0"; ValueType: string; ValueName: "ConfigFile"; ValueData: "{#ConfigINI}"
Root: HKLM64; Subkey: "Software\AdRem\GrafCrunch\1.0"; ValueType: string; ValueName: "DataFolder"; ValueData: "{#GrafCrunchProgramData}"

[INI]
Filename: {#ConfigINI}; Section: {#PathConfigSection}; Key: "data"; String: {#GrafCrunchProgramData}; Flags: createkeyifdoesntexist
Filename: {#ConfigINI}; Section: {#PathConfigSection}; Key: "logs"; String: {#GrafCrunchProgramData}\log; Flags: createkeyifdoesntexist

Filename: {#ConfigINI}; Section: {#GrafCrunchServerSection}; Key: "http_port"; String: "{code:GetGrafCrunchServerConfig|Port}"; Flags: createkeyifdoesntexist
Filename: {#ConfigINI}; Section: {#GrafCrunchServerSection}; Key: "domain"; String: "{code:GetGrafCrunchServerConfig|Domain}"; Flags: createkeyifdoesntexist

Filename: {#ConfigINI}; Section: {#NetCrunchServerConfigSection}; Key: "enable"; String: "true"; Flags: createkeyifdoesntexist
Filename: {#ConfigINI}; Section: {#NetCrunchServerConfigSection}; Key: "host"; String: "{code:GetNetCrunchServerConfig|Address}"; Flags: createkeyifdoesntexist
Filename: {#ConfigINI}; Section: {#NetCrunchServerConfigSection}; Key: "port"; String: "{code:GetNetCrunchServerConfig|Port}"; Flags: createkeyifdoesntexist
Filename: {#ConfigINI}; Section: {#NetCrunchServerConfigSection}; Key: "protocol"; String: "{code:GetNetCrunchServerConfig|Protocol}"; Flags: createkeyifdoesntexist
Filename: {#ConfigINI}; Section: {#NetCrunchServerConfigSection}; Key: "user"; String: "{code:GetNetCrunchServerConfig|User}"; Flags: createkeyifdoesntexist
Filename: {#ConfigINI}; Section: {#NetCrunchServerConfigSection}; Key: "password"; String: "{code:GetNetCrunchServerConfig|Password}"; Flags: createkeyifdoesntexist

[Dirs]
Name: {#GrafCrunchProgramData}

[Files]
Source: "SetupTools\Win32\Release\SetupTools.dll"; Flags: dontcopy 32bit

Source: {#MyAppIcon}; DestDir: "{app}"; Flags: ignoreversion
Source: {#LICENSE}; DestDir: "{app}"; Flags: ignoreversion
Source: {#NOTICE}; DestDir: "{app}"; Flags: ignoreversion

Source: "GrafCrunchGuard\Win64\Release\GrafCrunchGuard.exe"; DestDir: "{app}\bin\"; DestName: "GCGuard.exe"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "dest\bin\grafana-server.exe"; DestDir: "{app}\bin\"; DestName: "GCServer.exe"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "dest\conf\*"; DestDir: "{app}\conf\"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "dest\public\*"; DestDir: "{app}\public\"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "dest\vendor\*"; DestDir: "{app}\vendor\"; Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{group}\Uninstall"; Filename: "{uninstallexe}"; Comment: "Uninstall AdRem GrafCrunch Server";

[Run]
Filename: {app}\bin\GCGuard.exe; Parameters: "/install /silent"
Filename: {sys}\sc.exe; Parameters: "description GrafCrunchGuardService ""Provides infrastructure for AdRem GrafCrunch""" ; Flags: runhidden
FileName: {sys}\netsh; Parameters: "advfirewall firewall add rule name= ""AdRem GrafCrunch Server"" dir= in action= allow program= ""{app}\bin\GCServer.exe"" enable=yes"; Flags: runhidden
Filename: {sys}\sc.exe; Parameters: "start GrafCrunchGuardService" ; Flags: runhidden

[UninstallRun]
Filename: {sys}\sc.exe; Parameters: "stop GrafCrunchGuardService" ; Flags: runhidden
Filename: {app}\bin\GCGuard.exe; Parameters: "/uninstall /silent"
FileName: {sys}\netsh; Parameters: "advfirewall firewall delete rule name= ""AdRem GrafCrunch Server"""; Flags: runhidden

[UninstallDelete]
Type: files; Name: {#ConfigINI}

[Code]

type
  TNetCrunchServerConfig = record
    ErrorCode : Integer;
    Address : String;
    Port : String;
    Protocol : String;
    Password : String;
    Version : String;
  end;

function CompareVersion (AVer1, AVer2 : PAnsiChar) : Integer; external 'CompareVersion@files:SetupTools.dll stdcall';
function GetHostName : PAnsiChar; external 'GetHostName@files:SetupTools.dll stdcall';
function CheckServerPort (APort : PAnsiChar) : Integer; external 'CheckServerPort@files:SetupTools.dll stdcall';
function CheckNetCrunchWebAppServerConnection (AServerURL, User, Password: PAnsiChar) : Integer; external 'CheckNetCrunchWebAppServerConnection@files:SetupTools.dll stdcall';
function ReadNetCrunchServerConfig(AAddress, APort, APassword: PAnsiChar) : PAnsiChar; external 'ReadNetCrunchServerConfig@files:SetupTools.dll stdcall';

var
  UpdateGrafCrunchServer: Boolean;
  HostName: String;
  NetCrunchServerConfigData: TNetCrunchServerConfig;
  GrafCrunchServerConfig: TInputQueryWizardPage;
  NetCrunchServerConfig: TWizardPage;
  NetCrunchServerAddressTextBox: TEdit;
  NetCrunchServerPortTextBox: TEdit;
  NetCrunchServerPasswordTextBox: TPasswordEdit;
  InfoPage: TOutputMsgMemoWizardPage;

function CheckVersion : Boolean;
var 
  OldVersion: String;
  CurrentVersion: String;
  Uninstaller: String;
  ErrorCode : Integer;
begin
  CurrentVersion := '{#MyAppVersion}';
  if RegKeyExists(HKEY_LOCAL_MACHINE, '{#MyAppUninstallKey}') then begin
    RegQueryStringValue(HKEY_LOCAL_MACHINE, '{#MyAppUninstallKey}', 'DisplayVersion', OldVersion);

    if (CompareVersion(OldVersion, CurrentVersion) <= 0) then begin
      if (MsgBox('Version ' + OldVersion + ' of {#MyAppName} is already installed. Continue to use this old version?', mbConfirmation, MB_YESNO) = IDYES) then begin
        Result := False;
      end else begin
        RegQueryStringValue(HKEY_LOCAL_MACHINE, '{#MyAppUninstallKey}', 'UninstallString', Uninstaller);
        ShellExec('runas', Uninstaller, '/SILENT', '', SW_HIDE, ewWaitUntilTerminated, ErrorCode);
        Result := True;
      end;
    end else begin
      MsgBox('Version ' + OldVersion + ' of {#MyAppName} is already installed. This installer will exit.', mbInformation, MB_OK);
      Result := False;
    end;
  end else begin
    Result := True;
  end;
end;

function CreateLabel (AParent : TWizardPage; ALeft, ATop: Integer; const ACaption: String) : TLabel;
var TextLabel : TLabel;
begin
  TextLabel := TLabel.Create(AParent);
  with TextLabel do begin
    Parent := AParent.Surface;
    Left := ScaleX(ALeft);
    Top := ScaleY(ATop);
    Height := ScaleY(17);
    Caption := ACaption;
  end;
  Result := TextLabel;
end;

function CreateTextBox (AParent : TWizardPage; ATop, AWidth, ATabOrder : Integer; const ACaption: String) : TEdit;
var TextBox : TEdit;
begin
  CreateLabel(AParent, 0, ATop, ACaption);
  TextBox := TEdit.Create(AParent);
  with TextBox do begin
    Parent := AParent.Surface;
    Left := ScaleX(0);
    Top := ScaleY(ATop + 16);
    Width := ScaleX(AWidth);
    Height := ScaleY(25);
    TabOrder := ATabOrder;
    Text := '';
  end;
  Result := TextBox;
end;

function CreatePasswordTextBox (AParent : TWizardPage; ATop, AWidth, ATabOrder : Integer; const ACaption : String) : TPasswordEdit;
var PasswordTextBox : TPasswordEdit;
begin
  CreateLabel(AParent, 0, ATop, ACaption);
  PasswordTextBox := TPasswordEdit.Create(AParent);
  with PasswordTextBox do begin
    Parent := AParent.Surface;
    Left := ScaleX(0);
    Top := ScaleY(ATop + 16);
    Width := ScaleX(AWidth);
    Height := ScaleY(25);
    TabOrder := ATabOrder;
    Text := '';
  end;
  Result := PasswordTextBox;
end;

function CreateCheckBox (AParent : TWizardPage; ATop, ATabOrder : Integer; const ACaption: String) : TCheckBox;
var CheckBox : TCheckBox;
begin
  CheckBox := TCheckBox.Create(AParent);
  with CheckBox do begin
    Parent := AParent.Surface;
    Left := ScaleX(0);
    Top := ScaleY(ATop);
    Height := ScaleY(17);
    Width := AParent.SurfaceWidth;
    Caption := ACaption;
    Checked := False;
    TabOrder := ATabOrder;
  end;
  Result := CheckBox;
end;

function GrafCrunchServerDatabaseExist : Boolean;
begin
  Result := FileExists(ExpandConstant('{#GrafCrunchProgramData}' + '\grafana.db'));
end;

function GetDefaultData (const Section, Key, DefaultValue, PreviousDataKey: String) : String;
var
  NO_PREVIOUS_DATA_VALUE : String;
  DefaultData : String;

begin
  NO_PREVIOUS_DATA_VALUE := '#@$';
  DefaultData := GetPreviousData(PreviousDataKey, NO_PREVIOUS_DATA_VALUE);
  if (DefaultData = NO_PREVIOUS_DATA_VALUE) then begin
    DefaultData := GetIniString(ExpandConstant(Section), Key, ExpandConstant(DefaultValue), WizardDirValue + ExpandConstant('{#ConfigINIFile}'));
  end;
  Result := DefaultData;
end;

function GetDefaultGrafCrunchServerDomain : String;
var 
  DefaultServerDomain : String;

begin
  if (HostName <> '') 
    then DefaultServerDomain := HostName
    else DefaultServerDomain := '{#DefaultGrafCrunchServerDomain}';
  Result := GetDefaultData('{#GrafCrunchServerSection}', 'domain', DefaultServerDomain, 'GrafCrunchDomain');
end;

function GetDefaultGrafCrunchServerPort : String;
begin
  Result := GetDefaultData('{#GrafCrunchServerSection}', 'http_port', '{#DefaultGrafCrunchServerPort}', 'GrafCrunchPort');
end;

procedure SetGrafCrunchServerConfigDefaultValues;
begin
  GrafCrunchServerConfig.Values[0] := GetDefaultGrafCrunchServerDomain;
  GrafCrunchServerConfig.Values[1] := GetDefaultGrafCrunchServerPort;
end;

procedure PrepareGrafCrunchServerConfigPage;
begin
  GrafCrunchServerConfig := CreateInputQueryPage(wpSelectDir, 'GrafCrunch Server Configuration', '', 'Please specify GrafCrunch server settings, then click Next.');
  GrafCrunchServerConfig.Add('GrafCrunch server domain:', False);
  GrafCrunchServerConfig.Add('GrafCrunch server port:', False);
  SetGrafCrunchServerConfigDefaultValues;
end;

function GetGrafCrunchServerConfig(Param: String) : String;
begin
  Result := '';
  if (Param = 'Domain') then begin
    Result := GrafCrunchServerConfig.Values[0];
  end;
  if (Param = 'Port') then begin
    Result := GrafCrunchServerConfig.Values[1];
  end;
end;

function CheckGrafCrunchServerConfig : Boolean;
var
  Port : String;
  CheckResult : Integer;
  ErrorMessage : String;

begin
  Port := GetGrafCrunchServerConfig('Port');
  CheckResult := CheckServerPort(Port);
  if (CheckResult > 0) then begin
    case CheckResult of
      1: ErrorMessage := 'Port must be an integer value.';
      2: ErrorMessage := 'Port ' + Port + ' is reserved.';
      3:  ErrorMessage := 'Port ' + Port + ' is not available.';
      4:  ErrorMessage := 'Port must be greater than zero.';
    end;

    MsgBox(ErrorMessage, mbError, MB_OK);
    Result := False;
  end else begin
    Result := True;
  end;
end;

function GetDefaultNetCrunchServerAddress : String;
var 
  DefaultServerAddress : String;
begin
  if (HostName <> '') 
    then DefaultServerAddress := HostName
    else DefaultServerAddress := '{#DefaultNetCrunchServerAddress}';
  Result := GetDefaultData('{#NetCrunchServerConfigSection}', 'host', DefaultServerAddress, 'NetCrunchAddress');
end;

function GetDefaultNetCrunchServerPort : String;
var Port : String;
    DefaultNetCrunchPort : String;
begin
  Port := '';
  DefaultNetCrunchPort := '{#DefaultNetCrunchServerPort}'; 
  if not RegQueryStringValue(HKLM64, ExpandConstant('{#NetCrunchServerKey}'), 'Port', Port) then begin
    Port := '';
  end;
  if (Port <> '') then begin
    Result := Port;
  end else begin
    Result := GetDefaultData('{#NetCrunchServerConfigSection}', 'netcrunch_port', DefaultNetCrunchPort, 'NetCrunchPort');
  end;
end;

function GetDefaultNetCrunchServerPassword : String;
begin
  Result := '';
end;

procedure SetNetCrunchServerConfigDefaultValues;
begin
  NetCrunchServerAddressTextBox.Text := GetDefaultNetCrunchServerAddress;
  NetCrunchServerPortTextBox.Text := GetDefaultNetCrunchServerPort;
  NetCrunchServerPasswordTextBox.Text := GetDefaultNetCrunchServerPassword;
end;

procedure PrepareNetCrunchServerConfigPage;
begin
  NetCrunchServerConfig := CreateCustomPage(GrafCrunchServerConfig.ID, 'NetCrunch Server Configuration', '');
  CreateLabel(NetCrunchServerConfig, 0, 0, 'Please specify NetCrunch server settings, then click Next.');
  NetCrunchServerAddressTextBox := CreateTextBox (NetCrunchServerConfig, 24, NetCrunchServerConfig.SurfaceWidth, 0, 'NetCrunch server address:');
  NetCrunchServerPortTextBox := CreateTextBox (NetCrunchServerConfig, 76, NetCrunchServerConfig.SurfaceWidth, 1, 'NetCrunch server port:');
  NetCrunchServerPasswordTextBox := CreatePasswordTextBox (NetCrunchServerConfig, 128, NetCrunchServerConfig.SurfaceWidth, 2, 'NetCrunch server password:');
  SetNetCrunchServerConfigDefaultValues;
end;

function GetNetCrunchServerConnectionConfig(Param: String) : String;
begin
  Result := '';
  case Param of
    'Address': Result := NetCrunchServerAddressTextBox.Text;
    'Port': Result := NetCrunchServerPortTextBox.Text;
    'Password': Result := NetCrunchServerPasswordTextBox.Text;
  end;
end;

function DecodeServerConfig (EncodedServerConfig : String) : TNetCrunchServerConfig;
var
  ConfigList : TStringList;
  Config : TNetCrunchServerConfig;
  Position : Integer;
begin
  if (EncodedServerConfig <> '') then begin
    Config.ErrorCode := 12;
    ConfigList := TStringList.Create;
    try

      Position := 0;
      EncodedServerConfig := EncodedServerConfig + '~';
      while (EncodedServerConfig <> '') do begin
        Position := Pos('~', EncodedServerConfig);
        ConfigList.Add(Copy(EncodedServerConfig, 1, Position - 1));
        Delete(EncodedServerConfig, 1, Position);
      end;

      if (ConfigList.Count > 0) then begin
        Config.ErrorCode := StrToInt(ConfigList[0]);
        case Config.ErrorCode of
          0 : begin
                if (ConfigList.Count = 6) then begin
                  Config.Address := ConfigList[1];
                  Config.Version := ConfigList[2];
                  Config.Password := ConfigList[3];
                  Config.Port := ConfigList[4];
                  Config.Protocol := ConfigList[5];
                end else begin
                  Config.ErrorCode := 12;
                end;
             end;
          10 : begin
                 if (ConfigList.Count = 2)
                   then Config.Version := ConfigList[1]
                   else Config.ErrorCode := 12;
               end;
        end;
      end;
    finally
      ConfigList.Free;
    end;
  end else begin
    Config.ErrorCode := 12;
  end;
  Result := Config;
end;

function CheckNetCrunchServerConfig : Boolean;
var 
  Address : String;
  Port : String;
  Password : String;
  EncodedServerConfig : String;
  ErrorMessage : String;
begin
  Address := GetNetCrunchServerConnectionConfig('Address');
  Port := GetNetCrunchServerConnectionConfig('Port');
  Password := GetNetCrunchServerConnectionConfig('Password');
  Result := False;

  if ((Address <> '') and (Port <> '')) then begin
    EncodedServerConfig := ReadNetCrunchServerConfig(Address, Port, Password);
    NetCrunchServerConfigData := DecodeServerConfig(EncodedServerConfig);

    if (NetCrunchServerConfigData.ErrorCode <> 0) then begin
      case NetCrunchServerConfigData.ErrorCode of
        -1 : ErrorMessage := 'Can' + #39 + 't connect to NetCrunch server';
         3 : ErrorMessage := 'Incorrect password';
         4 : ErrorMessage := 'Session limit reached';
         5 : ErrorMessage := 'Remote access disabled';
         6 : ErrorMessage := 'No remote access license';
         7, 8 : ErrorMessage := 'Unknown error';
         9 : ErrorMessage := 'GrafCrunch requires NetCrunch server version 9.0.0 or higher.';
        10 : ErrorMessage := 'GrafCrunch requires NetCrunch server version 9.0.0 or higher. Your version is ' + NetCrunchServerConfigData.Version;
        11 : ErrorMessage := 'Incorrect NetCrunch server address';
        12 : ErrorMessage := 'Conversion error';
      end;
      MsgBox(ErrorMessage, mbError, MB_OK);
    end else begin
      Result := True;
    end;
  end;
end;

function GetNetCrunchServerConfig(Param: String) : String;
begin
  Result := '';
  if (NetCrunchServerConfigData.ErrorCode = 0) then begin
    case Param of
      'Address': Result := NetCrunchServerConfigData.Address;
      'Port': Result := NetCrunchServerConfigData.Port;
      'Protocol': Result := NetCrunchServerConfigData.Protocol;
      'User': Result := '{#DefaultNetCrunchServerUser}';
      'Password': Result := NetCrunchServerConfigData.Password;
    end;
  end;
end;

procedure PrepareInfoPage;
begin
  InfoPage := CreateOutputMsgMemoPage(wpInstalling, 'GrafCrunch server connection info', '', '', '');
  with InfoPage.RichEditViewer do begin
    Font.Size := 8;
    Lines.Add('');
    Lines.Add('  GrafCrunch is available at: ' + 'http://' + GetGrafCrunchServerConfig('Domain') + ':' + GetGrafCrunchServerConfig('Port'));
    Lines.Add('');
  end;

  if UpdateGrafCrunchServer then begin
    with InfoPage.RichEditViewer do begin
      Lines.Add('  Previous GrafCrunch server data were detected.');
      Lines.Add('  You can log in as previously defined user.');
    end;
  end else begin
    with InfoPage.RichEditViewer do begin
      Lines.Add('  Default GrafCrunch admin was created with credentials:');
      Lines.Add('    User: admin');
      Lines.Add('    Password: admin');
    end;
  end;
end;

procedure InitializeWizard;
begin
  if not CheckVersion then Abort;
  UpdateGrafCrunchServer := GrafCrunchServerDatabaseExist;
  HostName := String(AnsiString(GetHostName));
  PrepareGrafCrunchServerConfigPage;
  PrepareNetCrunchServerConfigPage;
  PrepareInfoPage;
end;

procedure RegisterPreviousData(PreviousDataKey: Integer);
begin
  SetPreviousData(PreviousDataKey, 'GrafCrunchDomain', GetGrafCrunchServerConfig('Domain'));
  SetPreviousData(PreviousDataKey, 'GrafCrunchPort', GetGrafCrunchServerConfig('Port'));
  SetPreviousData(PreviousDataKey, 'NetCrunchAddress', GetNetCrunchServerConfig('Address'));
  SetPreviousData(PreviousDataKey, 'NetCrunchPort', GetNetCrunchServerConfig('Port'));
end;

function NextButtonClick(CurPageID: Integer): Boolean;
var CheckStatus : Boolean;
begin
  CheckStatus := True;
  case CurPageID of
    GrafCrunchServerConfig.ID: CheckStatus := CheckGrafCrunchServerConfig;
    NetCrunchServerConfig.ID: CheckStatus := CheckNetCrunchServerConfig;
  end;
  Result := CheckStatus;
end;

function UpdateReadyMemo(Space, NewLine, MemoUserInfoInfo, MemoDirInfo, MemoTypeInfo, MemoComponentsInfo, MemoGroupInfo, MemoTasksInfo: String): String;
begin

  Result := MemoDirInfo + NewLine + NewLine +
    'Data location:' + NewLine +
      Space + ExpandConstant('{#GrafCrunchProgramData}') + NewLine + NewLine + 
    
    'GrafCrunch server settings:' + NewLine +   
      Space + 'Domain: ' + GetGrafCrunchServerConfig('Domain') + NewLine + 
      Space + 'Port: ' + GetGrafCrunchServerConfig('Port') + NewLine + NewLine + 
          
    'NetCrunch web server settings:' + NewLine + 
      Space + 'Address: ' + GetNetCrunchServerConfig('Address') + NewLine +
      Space + 'Port: ' + GetNetCrunchServerConfig('Port') + NewLine + 
      Space + 'Protocol: ' + GetNetCrunchServerConfig('Protocol');
end;

//**************

//Kiedy NC nie jest zainstalowany na mszynie gdzie jest instalowany GC to domy�lny adres NC serwera powinien byc pusty

//Add proceses descriptions

//Add shortcuts for start / stop GrafCrunch service
//Name: "{group}\Start GrafCrunch Server"; Filename: {sys}\sc.exe; Parameters: "start GrafCrunchGuardService" ; Flags: runminimized; IconFilename: {app}\{#MyAppIcon}; Comment: "Starts AdRem GrafCrunch Server";
//Name: "{group}\Stop GrafCrunch Server"; Filename: {sys}\sc.exe; Parameters: "stop GrafCrunchGuardService" ; Flags: runminimized; IconFilename: {app}\{#MyAppIcon}; Comment: "Stops AdRem GrafCrunch Server";

//Implement Modify mode for server config modifications;
//Grafana server log problem

//;Filename: {#ConfigINI}; Section: {#NetCrunchServerConfigSection}; Key: "user"; String: {#NetCrunchServerUser}; Flags: createkeyifdoesntexist
//;Filename: {#ConfigINI}; Section: {#NetCrunchServerConfigSection}; Key: "password"; String: {#NetCrunchServerPassword}; Flags: createkeyifdoesntexistcls
