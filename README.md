# Dagger

Ce dépôt contient les ressources relatives à l'épisode 33 de inpulse.tv 👉 https://youtu.be/m5mauLiCp3Y

## Api playground 

Retrouvez le playground [ici](https://play.dagger.cloud/playground)

Et l'exemple plus fourni présenté dans la vidéo [ici](https://play.dagger.cloud/playground/Z_1Px4cAIsr)

## Go sdk

### Pré-requis 

🧇 [Go 1.15 ou ultérieur](https://go.dev/doc/install)

🧙 [Mage](https://magefile.org/)

🐋 [Docker](https://docs.docker.com/get-docker/)

###

Pour voir l'ensemble des targets mage que vous pouvez lancer dans le dossier racine
``` bash
mage 
```

## Pour aller plus loin 

### Github actions

A la fin de la vidéo, j'évoque l'intégration de notre CI en GO dans une plateforme comme github actions. 
Vous pouvez retrouver l'intégration minimaliste [ici](.github/workflows/ci.yaml)

> Dagger est en train de s'éméanciper un maximum des dépendances aux différentes plateformes de ci

### CI plus complexe

Au vu des capacités d'un langage comme go, la ci présentée reste simpliste. 

Vous pouvez vous rendre sur la branche dev de ce repo pour trouver un exemple plus complet faisant intervenir le capacité de programmation concurrente de go dans le cadre de compilation multi architecture.
