# gui

This is the Parallelcoin Pod wallet GUI, built on a fork of
[gio](https://gioui.org) and using a customised set of widgets and layouts from
in the layout, widget and widget/material folders of that repository, where
several additional widget types have been added, and mostly this is in sync
with the mainline of the gio repository, though ultimately either they will
be copying this, or gio will be completely forked as it uses a completely
different programming idiom that doesn't aim at arbitrary abstract rules
about composition of the code versus the state, because of the difficulty
that the pure immediate mode programming model requires. This is because
that idiom is extremely wordy and long winded and noisy and difficult to
read, rearrange, and fix.

It has a run control system to launch the individual modules of the
Parallelcoin Pod, including wallet, blockchain node, a CLI interface to
these two modules. 