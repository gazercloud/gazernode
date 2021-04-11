
!define ARP "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"
!include "MUI2.nsh"
!include "FileFunc.nsh"

!define MUI_HEADERIMAGE
;!define MUI_HEADERIMAGE_BITMAP "installer\installer-icon.bmp"

;!define MUI_COMPONENTSPAGE_SMALLDESC ;No value
;!define MUI_UI "myUI.exe" ;Value
;!define MUI_INSTFILESPAGE_COLORS "FFFFFF 000000" ;Two colors

Name "${PRODUCT_NAME}"
Caption "Install ${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "${PRODUCT_NAME}_${PRODUCT_VERSION}.exe"
ShowInstDetails show

InstallDir "$PROGRAMFILES64\gazer_node"

;Page directory
;Page instfiles

SetCompressor /SOLID lzma

!define MUI_ABORTWARNING

!insertmacro MUI_PAGE_LICENSE "license.txt"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

    # These indented statements modify settings for MUI_PAGE_FINISH
	!define MUI_FINISHPAGE_NOAUTOCLOSE
	!define MUI_FINISHPAGE_RUN
	!define MUI_FINISHPAGE_RUN_CHECKED
	!define MUI_FINISHPAGE_RUN_TEXT "Start Gazer Node"
	!define MUI_FINISHPAGE_RUN_FUNCTION "LaunchLink"
!insertmacro MUI_PAGE_FINISH
  
!insertmacro MUI_LANGUAGE "English"
  
Section
	SetOutPath "$PROGRAMFILES64\gazer_node"

	ExecWait '"$INSTDIR\gazer_node.exe" -stop'
	ExecWait '"$INSTDIR\gazer_node.exe" -uninstall'

	File "..\bin\gazer_node.exe"
	
	WriteUninstaller $INSTDIR\uninstaller.exe
	
	WriteRegStr HKLM "${ARP}" "DisplayName" "${PRODUCT_NAME}"
	WriteRegStr HKLM "${ARP}" "UninstallString" "$\"$INSTDIR\uninstaller.exe$\""
	WriteRegStr HKLM "${ARP}" "Publisher" "Gazer.Cloud"
	WriteRegStr HKLM "${ARP}" "DisplayVersion" "${PRODUCT_VERSION}"
	WriteRegStr HKLM "${ARP}" "DisplayIcon" "$INSTDIR\gazer_node.exe"
	
	
	${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
	IntFmt $0 "0x%08X" $0
	WriteRegDWORD HKLM "${ARP}" "EstimatedSize" "$0"
	
	
	CreateShortCut "$DESKTOP\${PRODUCT_NAME}.lnk" "$INSTDIR\gazer_node.exe" ""
	CreateDirectory "$SMPROGRAMS\${PRODUCT_NAME}"
	CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\uninstaller.lnk" "$INSTDIR\uninstaller.exe" "" "$INSTDIR\uninstaller.exe" 0
	CreateShortCut "$SMPROGRAMS\${PRODUCT_NAME}\${PRODUCT_NAME}.lnk" "$INSTDIR\gazer_node.exe" "" "$INSTDIR\gazer_node.exe" 0

	ExecWait '"$INSTDIR\gazer_node.exe" -install'
	ExecWait '"$INSTDIR\gazer_node.exe" -start'
SectionEnd

Section Uninstall
	ExecWait '"$INSTDIR\gazer_node.exe" -stop'
	ExecWait '"$INSTDIR\gazer_node.exe" -uninstall'

	Delete $INSTDIR\uninstaller.exe
	Delete $INSTDIR\gazer_node.exe
	
	Delete "$DESKTOP\${PRODUCT_NAME}.lnk"
	Delete "$SMPROGRAMS\${PRODUCT_NAME}\*.*"
	RmDir  "$SMPROGRAMS\${PRODUCT_NAME}"	
	
	DeleteRegKey HKLM "${ARP}"
	
	RMDir $INSTDIR
SectionEnd

Function LaunchLink
  ExecShell "" "$INSTDIR\gazer_node.exe"
FunctionEnd
