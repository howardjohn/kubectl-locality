# Kubectl Locality Plugin

A plugin to get the locality of pods

## Install

`go install github.com/howardjohn/kubectl-locality@latest`

## Usage

Example output:

```
$ kubectl locality -n default
NAMESPACE     NAME                      REGION          ZONE
default       echo-cb96f8d94-2j6c6      us-central1     us-central1-f
default       echo-cb96f8d94-5kcgr      us-central1     us-central1-f
default       echo-cb96f8d94-86pnt      us-central1     us-central1-a
default       echo-cb96f8d94-jjr6r      us-central1     us-central1-b
default       echo-cb96f8d94-p77bk      us-central1     us-central1-b
default       echo-cb96f8d94-t9b58      us-central1     us-central1-a
default       echo-cb96f8d94-wgvlv      us-central1     us-central1-a
default       echo-cb96f8d94-zvqtd      us-central1     us-central1-b
default       shell-7854df9c5-ggclf     us-central1     us-central1-a
```
