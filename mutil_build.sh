#!/usr/bin/env bash

buildPath="build"
repoPath="repo"
ReNameCode="qiniu_upload"
VERSION_MAJOR=1
VERSION_MINOR=0
VERSION_PATCH=0
VERSION_BUILD=0

VersionCode=$[$[VERSION_MAJOR * 100000000] + $[VERSION_MINOR * 100000] + $[VERSION_PATCH * 100] + $[VERSION_BUILD]]
VersionName="${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}.${VERSION_BUILD}"
packageReName="${ReNameCode}_${VersionName}"

shell_running_path=$(cd `dirname $0`; pwd)

checkFuncBack(){
  if [ $? -eq 0 ]; then
    echo -e "\033[;30mRun [ $1 ] success\033[0m"
  else
    echo -e "\033[;31mRun [ $1 ] error exit code 1\033[0m"
    exit 1
  fi
}

checkEnv(){
  evn_checker=`which $1`
  checkFuncBack "which $1"
  if [ ! -n "evn_checker" ]; then
    echo -e "\033[;31mCheck event [ $1 ] error exit\033[0m"
    exit 1
  else
    echo -e "\033[;32mCli [ $1 ] event check success\033[0m\n-> \033[;34m$1 at Path: ${evn_checker}\033[0m"
  fi
}

if [ -d "${buildPath}" ]; then
    rm -rf ${buildPath}
    sleep 1
fi

checkEnv tar

echo -e "============\nPrint build info start"
go version
which go
echo -e "Your settings is
\tVersion Name -> ${ReNameCode}
\tVersion code -> ${VersionCode}
\tVersion name -> ${VersionName}
\tPackage rename -> ${packageReName}
\tOut Path -> ${shell_running_path}/${buildPath}
"
echo -e "Print build info end\n============"

mkdir -p ${buildPath}
echo "start build OSX 64"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
mv main "${buildPath}/${packageReName}_osx_64"
echo "build OSX 64 finish"

echo "start build OSX 32"
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build main.go
mv main "${buildPath}/${packageReName}_osx_86"
echo "build OSX 32 finish"

echo "start build Linux 64"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
mv main "${buildPath}/${packageReName}_linux_64"
echo "build linux 64 finish"

echo "start build Linux 32"
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build main.go
mv main "${buildPath}/${packageReName}_linux_86"
echo "build linux 32 finish"

echo "start build windows 64"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
mv main.exe "${buildPath}/${packageReName}_win_64.exe"
echo "build windows 64 finish"

echo "start build windows 32"
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build main.go
mv main.exe "${buildPath}/${packageReName}_win_86.exe"
echo "build windows 32 finish"

read -p "Do you want repo:(nothing is not need)? " word
if [ -n "$word" ] ;then
    out_folder=${VERSION_MAJOR}-${VERSION_MINOR}-${VERSION_PATCH}
    out_tar_name="${ReNameCode}-${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}.tag.gz"
    out_folder_path="${repoPath}/${ReNameCode}/${out_folder}/"
    out_tar_path="${repoPath}/${ReNameCode}/${out_tar_name}"
    mkdir -p "${out_folder_path}"
    cp -r "${buildPath}/" "${out_folder_path}"
    echo -e "Repo config
\tRepo path: ${out_folder_path}
\tRepo out_tar_path: ${out_tar_path}
"
    cat > "${out_folder_path}Watch.bat" << EOF
@echo off
@echo. ==== start watch info ====
${packageReName}_win_86.exe -s "%~nx1"
pause
EOF
    cat > "${out_folder_path}Watch_64.bat" << EOF
@echo off
@echo. ==== start watch info ====
${packageReName}_win_64.exe -s "%~nx1"
pause
EOF
    cat > "${out_folder_path}MD5.bat" << EOF
@echo off
@echo. ==== start watch info ====
${packageReName}_win_86.exe -m "%~nx1"
pause
EOF
    cat > "${out_folder_path}MD5_64.bat" << EOF
@echo off
@echo. ==== start watch info ====
${packageReName}_win_64.exe -m "%~nx1"
pause
EOF

cd ${repoPath}
tar zcvf "${ReNameCode}/${out_tar_name}" "${ReNameCode}/${out_folder}"
cd ${shell_running_path}
rm -rf ${out_folder_path}
sleep 1

echo -e "=>Out repo tar \033[;36m-> ${out_tar_path}\033[0m"
fi

echo -e "============\nAll the build is finish! at Build Path\n${shell_running_path}/${buildPath}"
