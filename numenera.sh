#! /bin/sh

ALLCONFS="./conf/numenera.yaml \
./conf/character.yaml \
./conf/completion.yaml \
./conf/danger.yaml \
./conf/dice.yaml \
./conf/location.yaml \
./conf/name.yaml \
./conf/util.yaml \
./conf/player.yaml"

HELPSTR="numenera.sh [num] [category]

num =  number of generations

category =
            'character' will generate a character idea
            'location' will generate a location idea
            'danger' will generate a danger idea
            'campaign' will generate a campaign idea, default"

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
elif [ "danger" = "$2" ]; then
  INITCONF="./conf/danger.yaml:"
elif [ "campaign" = "$2" ]; then
  INITCONF=""
fi

CONFFLAGS=""
if [ -n "$INITCONF" ]; then
  CONFFLAGS="-sb ${INITCONF}"
fi
for CONF in $ALLCONFS; do
  CONFFLAGS="${CONFFLAGS} -sb ${CONF}"
done

SEED="$3"
if [ -z $SEED ]; then 
  SEED=$(cat /dev/urandom |  tr -dc '[:alpha:]' | fold -w ${1:-20} | head -n 1 )
  echo "Seed: $SEED"
fi

if [ -n "$1" ]; then
    buildstory -r $1 -seed $SEED $CONFFLAGS
else
    buildstory -seed $SEED $CONFFLAGS
fi
