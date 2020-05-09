@echo off
@ cls
echo "**Mensajeria P2p Cifrada**"
echo .
echo Levanta: ipfs daemon 

:inicio

set /p texto="Usted:> " 
msg.exe -e --llave hola --txt %texto%
:set /p cont="secuencia:> "
:ipfs name publish %cont%
msg.exe -d --llave hola --txt .
goto inicio