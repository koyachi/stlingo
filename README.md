# stlingo

Collect similar strings in Android's strings.xml and iOS's Localizable.strings files.

## Install

 $ go install github.com/koyachi/stlingo

## Usage

  $ stlingo /path/to/in.txt > /path/to/out.csv

## in.txt format

  // comment
  ./example/android_strings.xml
  ./example/ios_Localizable.strings

## example output

  >> read file: ./example/input.txt
  ,index, platformType, hash, filePath
  , 0, Android, 5b3e6697bd23c3bdfd1c3dbcfa3e3000, /Users/koyachi/src/github.com/koyachi/stlingo/example/android_strings.xml
  , 1, IOS, 13546c3bf12e359d66570d09ced7826b, /Users/koyachi/src/github.com/koyachi/stlingo/example/ios_Localizable.strings
  
  
  ,PlatformType, score(or X), val, key, line, file_index
  ,Android,X,"bar",foo,0,/Users/koyachi/src/github.com/koyachi/stlingo/example/android_strings.xml
  ,IOS,0,"bar",foo,1,1
  ,Android,1,"baz",foo1,0,0
  ,IOS,3,"入力",ios_label_2,3,1
  ,Android,3,"入力",android_label_2,0,0
  ,Android,X,"baz",foo1,0,/Users/koyachi/src/github.com/koyachi/stlingo/example/android_strings.xml
  ,IOS,3,"入力",ios_label_2,3,1
  ,Android,3,"入力",android_label_2,0,0
  ,IOS,X,"入力",ios_label_2,3,/Users/koyachi/src/github.com/koyachi/stlingo/example/ios_Localizable.strings
  ,Android,0,"入力",android_label_2,0,0
  ,IOS,3,"入力エラー",ios_label_1,2,1
  ,Android,3,"入力エラー",android_label_1,0,0
  ,IOS,3,"文字を入力",ios_label_3,4,1
  ,Android,3,"文字を入力",android_label_3.0,0,0
  ,IOS,X,"入力エラー",ios_label_1,2,/Users/koyachi/src/github.com/koyachi/stlingo/example/ios_Localizable.strings
  ,Android,0,"入力エラー",android_label_1,0,0
  ,IOS,X,"文字を入力",ios_label_3,4,/Users/koyachi/src/github.com/koyachi/stlingo/example/ios_Localizable.strings
  ,Android,0,"文字を入力",android_label_3.0,0,0
  ,Android,2,"文字を削除",android_label_3.1,0,0
  ,Android,X,"文字を削除",android_label_3.1,0,/Users/koyachi/src/github.com/koyachi/stlingo/example/android_strings.xml

