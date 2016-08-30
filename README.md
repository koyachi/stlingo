# stlingo

Collect similar strings in Android's strings.xml and iOS's Localizable.strings files.

## Install

 $ go install github.com/koyachi/stlingo

## Usage

  $ stlingo /path/to/in.txt > /path/to/out.csv

## in.txt format

  // comment: iOS
  /path/to/ios_app/Localizable.strings
  /path/to/ios_app/Localizable2.strings
  
  // comment: Android
  /path/to/android_app/src/main/res/values/strings.xml

