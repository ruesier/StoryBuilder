#! /bin/sh

if [ "help" = $1 ]; then
  echo "numenera.sh [num] [category]

num =  number of generations

category =
            'character' will generate a character idea
            'location' will generate a location idea
            otherwise, will generate a campaign idea
  "
return 0
fi

ALLCONFS="./conf/numenera.yaml:./conf/character.yaml:./conf/completion.yaml:./conf/danger.yaml:./conf/dice.yaml:./conf/location.yaml:./conf/name.yaml:./conf/util.yaml"
INITCONF=""

if [ "character" = "$2" ]; then
  INITCONF="./conf/character.yaml:"
elif [ "location" = "$2" ]; then
  INITCONF="./conf/location.yaml:"
fi

buildstory -sb "$INITCONF$ALLCONFS" -r $1
