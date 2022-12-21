# The ParallelCoin Pod

This is an all-in one, single binary, monorepo that contains the blockchain node, wallet, miner, CLI RPC client, and GUI for both the current legacy chain, and a substantial upgrade to the protocol that will fork the chain in the near future to bring the network up to date and fix its difficulty adjustment problems, and introduce a new proof of work and multi-interval block schedule that improves the chain's precision in difficulty adjustment.

## Building

Firstly, you will need to install prerequisites.

For the majority of users, this will mean ubuntu/debian based linux, for which you can get set up correctly with the script [prereqs/ubuntu.sh](prereqs/ubuntu.sh)

Note that some of the generators for the GUI library's shaders will not work, but the code is already working that you can find in [shaders](shaders/). There is windows programs for this, they will be updated at some stage in the future based on the Gio library that we integrated here.

The versions for Fedora and Arch Linux can be found in [prereqs/](prereqs/) also. Note that the Fedora script has not been tested in a long time. The Arch version is the most up to date, as Arch is the project's preferred and recommended linux distribution (or manjaro, second best). 

Arch linux has the latest version available, if you are using an AUR helper like `yay` to install packages, but for everyone else, you can use [prereqs/go.sh](prereqs/go.sh) to automatically install Go for your system.

Next, run the script [build.sh](build.sh) and this will install `pod` in `$HOME/bin/`

> WARNING: the software is not stable or fully functional yet

After you have run `build.sh` from then on you can use `buidl install` to update it from the source code, and update `buidl` with `buidl builder`.

## Why Parallelcoin?

Parallelcoin was a fork of bitcoin that was released by a hit-and-run, cheap and dirty fork of bitcoin in early 2014 by a [bitcointalk.org](https://bitcointalk.org) user called "paralaxis". It merged the proof of work used with Litecoin and created a poorly thought out combined block timing schedule that has over time slowly decelerated due to its hard adjustment limiter and periodic bouts of cloud mining jumping all over it when difficulty finally comes down after the previous hit.

The code in this repository was primarily created by David Vennik, with some of the Gio custom widgets created by Djordje Marcetin (marcetin on Bitcointalk). David is the main designer of the current existing codebase, and wrote a very large amount of supporting libraries as part of this process, from configuration, to merging the btcd and btcwallet into one, adapting the protocol to work with Parallelcoin, and then designing the hard fork system, created a multicast mining cluster control system, CPU based proof of work and the bulk of the current GUI.

It is once again the bear market, as it was in the time that this project was started. This time, David has some savings and is working to get this to beta standard and promoting it to the #bitcoin twitter. David has been paid very little to nothing to do all this work, and is aiming to keep the ethical and philosophical purity of this project to enable it to be the first non-security cryptocurrency project. As such there is no on chain subsidies, no governance system, and no plans to ever allow such a thing.

Future plans include integrating the latest code from btcd and btcwallet to add Segwit support and to integrate and fork Lightning Network's LND and fork Neutrino to work with the chain, and ultimately to enable cross chain atomic swaps with Bitcoin, creating a pure decentralised marketplace and in-bitcoin valuation for Parallelcoin.
