# proc

This is a unified library for process control and logging.

Originally this was 4 separate libraries, and the tangle of overlap between them became a source of bugs.

What is provided in here is a reliable way to start new processes and stop them, and pipe their log entries out to
the controlling TTY, and receive and correctly process signals from outside a process most specifically for 
interrupt, in a fully uniform manner, that can be plugged into any application and make it usable as a worker for 
another application.

In future this library will add integrated Windows server process control interfacing, as part of the reason for the 
mess was precisely handling the differences between unix process signals and the windows process control libraries, 
which are less advanced.