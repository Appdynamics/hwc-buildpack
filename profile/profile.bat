:: This file's contents to be executed as a bat file on startup.
@echo off
powershell.exe -ExecutionPolicy Unrestricted %~dp0\profile.ps1 %1
