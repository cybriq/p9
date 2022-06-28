# jrnl

A simple go script that automatically generates a file based on unix
datestamp that I use to write journal entries in my private journal, as I am
a prolific writer with zero audience. It expects a JSON formatted
configuration file in the user's home directory called `.jrnl` to specify
the filesystem path where you want to keep these datestamped files for a
journal, with one element, a string with the key "Root", eg:

```json
{
  "Root": "path/to/journal/folder"
}
```