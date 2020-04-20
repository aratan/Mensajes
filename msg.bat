@echo off
@ cls
echo "**Mensajeria P2p Cifrada**"
echo Levanta: btfs daemon 

:inicio
set /p texto="texto a cifrar:> " 
btfs-codificadorDes -e --llave hola --txt %texto%

set /p publica="Clave de publicacion: "
btfs name publish %publica%
btfs-codificadorDes -d --llave hola --txt .
pause
goto inicio