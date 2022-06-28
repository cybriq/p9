# swd

Save Working Directory

## What is this?

For some reason, bash shell makes it incredibly difficult to put a function
into the prompt script that saves the current working directory, so that
when a session is closed, you go back to the same folder.

This is a little thing written in Go that does this without fuss, directly,
meaning I don't have to type the long path to my last working directory
every time I restart my machine or close a terminal I was using.