{****************************************************************
 *
 * Author : Boguslaw Gorczyca
 * Created: 2015-10-06
 *
 * No part of this file may be duplicated, revised, translated,
 * localized or modified in any manner or compiled, linked or
 * uploaded or downloaded to or from any computer system without
 * the prior written consent of AdRem Software sp z o.o.
 *
 * 2010 Copyright AdRem Software, all rights reserved
 ****************************************************************}

program GrafCrunchGuard;
uses
  SysUtils,
  Vcl.SvcMgr,
  madExcept,
  madLinkDisAsm,
  madListHardware,
  madListProcesses,
  madListModules,
  uMessages,
  uConfigFilesUtils,
  uNetCrunchConsts,
  uErrorReportingUtils,
  uGrafCrunchGuardService in 'uGrafCrunchGuardService.pas' {GrafCrunchGuardService: TService},
  uServiceDebugForm in 'V:\ADREM\Debug\uServiceDebugForm.pas' {ServiceDebugForm};

{$R *.RES}

const
  LOGS_FOLDER = 'C:\ProgramData\AdRem\GrafCrunch\log\';

begin
  ForceDirectories(LOGS_FOLDER);
  MessagesList.LogLevel := tmsMessage;
  MessagesList.LogFile := LOGS_FOLDER + ChangeFileExt(ExtractFileName(ParamStr(0)), '.log');
  MessagesList.WriteToFile := True;
  MessagesList.ReduceLogFileSize(500000);

  if not Application.DelayInitialize or Application.Installing then Application.Initialize;
  Application.CreateForm(TGrafCrunchGuardService, GrafCrunchGuardService);

  // For debugging purpose
  // Application.CreateForm(TServiceDebugForm, ServiceDebugForm);

  InitializeErrorReporting(True);
  Application.Run;
end.
