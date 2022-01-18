#! /bin/sh

ALLCONFS="./conf/numenera.yaml:./conf/character.yaml:./conf/completion.yaml:./conf/danger.yaml:./conf/dice.yaml:./conf/location.yaml:./conf/name.yaml:./conf/util.yaml"

HELPSTR="numenera.sh [num] [category]

num =  number of generations

category =
            'character' will generate a character idea
            'location' will generate a location idea
            otherwise, will generate a campaign idea"

if [ -n "$1" ]; then
    if [ "help" = $1 ]; then
        echo $HELPSTR
        return 0
    fi
fi

INITCONF=""

if [ "character" = "$2" ]; then
  INITCONF="./conf/character.yaml:"
elif [ "location" = "$2" ]; then
  INITCONF="./conf/location.yaml:"
fi

if [ -n "$1" ]; then
    buildstory -sb "$INITCONF$ALLCONFS" -r $1
else
    buildstory -sb "$ALLCONFS"
fi
