unset CDPATH
script="`test -L "$0" && readlink -n "$0" || echo "$0"`"
dir="$PWD"
cd "`dirname "$script"`"

if [ `echo "$1" | sed -e "s/\(.\).*/\1/"` = / ]
	then
		arg1="$1"
	else
		arg1="${dir}/$1"
	fi

echo Assembling "$arg1"
java -classpath "${CLASSPATH}:bin/classes:bin/lib/Hack.jar:bin/lib/HackGUI.jar:bin/lib/Compilers.jar:bin/lib/AssemblerGUI.jar:bin/lib/TranslatorsGUI.jar" HackAssemblerMain "$arg1"