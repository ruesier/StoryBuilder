#! /bin/sh


ALLCONFS="./conf/numenera.yaml:./conf/character.yaml:./conf/completion.yaml:./conf/danger.yaml:./conf/dice.yaml:./conf/location.yaml:./conf/name.yaml:./conf/util.yaml"

buildstory -sb "$ALLCONFS"