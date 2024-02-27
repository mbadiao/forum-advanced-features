#!/bin/bash
git config --global credential.helper store
git config --global user.email "elie.jean.malo@gmail.com"
git config --global user.name "emalo"

echo -n "enter the message: "
read var
git add .
git commit -m "$var"
git push