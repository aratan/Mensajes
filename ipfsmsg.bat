@echo off
@ cls
echo "**Mensajeria P2p Cifrada**"
echo Levanta: ipfs daemon 

:inicio
set /p texto="Usted:> " 
ipfsmsg.exe -e --llave hola --txt %texto%
set /p cont="secuencia:> "
ipfs name publish %cont%
ipfsmsg.exe -d --llave hola --txt .
goto inicio