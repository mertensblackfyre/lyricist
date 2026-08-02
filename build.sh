#!/bin/bash

make build
sudo rm /usr/sbin/lyricist
sudo mv lyricist /usr/sbin/lyricist
