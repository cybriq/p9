# The ParallelCoin Pod

This is an all-in one, single binary, monorepo that contains the blockchain 
node, wallet, miner, CLI RPC client, and GUI for both the current legacy 
chain, and a substantial upgrade to the protocol that will fork the chain in 
the near future to bring the network up to date and fix its difficulty 
adjustment problems, and introduce a new proof of work and multi-interval 
block schedule that improves the chain's precision in difficulty adjustment.

## Building

Please note that the following instructions have been tested and worked on 
Pop OS 21 (should be the same for ubuntu 21, also should work on ubuntu 20). 
The generators may not work but the prereqs have been updated to add things 
that seem to be missing for ubuntu 21 base.

The ParallelCoin Pod is by default dependent on go 1.16 or later. It likely
can build on earlier versions but newer is generally better, at least with
the Go compiler codebase, unlike with many other languages like Java and C++
where it's a crap shoot whether it's a good idea to upgrade (eg, Goland's
linux implementation since the beginning of 2021).

In the `prereqs/` folder you can find scripts that will install the needed
dependencies on Ubuntu, Fedora and Arch Linux. (apt, rpm and yay,
respectively)

To build correctly, so the versions are updated, first you should build the
builder script:

```bash
cd path/to/repository/root
go install ./pod/podbuild/.
```

This assumes you have correctly put the $GOBIN environment variable path
into your shell's path, as `go build` has the undesirable behaviour of
dropping the binary into the repository filesystem tree. The same applies
for the rest of these instructions.

Next, to update all the things (for now just ignore the generator failures
relating to the gio shader compilation, they are pre-generated):

```bash
podbuild generate
```

Next, build and install the binary into your $GOBIN folder:

```bash
podbuild install
```

And now, `pod` will get you started and show you the wallet creation screen,
unless you already did that.

If you are building for a headless, server version, such as a chain/wallet
for supporting applications using blockchain/wallet data, use the tag
`headless` to leave out the GUI component.

### Other platforms

The instructions above basically are the same, except for the differences in
how to set up environment variables, on Mac, Windows, and FreeBSD.

To build the android version, you can use the gogio fork found at
`pkg/gel/gio/gogio`. The same tooling should allow also building for iOS
devices.

There is a `appstore` build tag that removes the miner code from
the output binaries, as the duopoly of app curation have agreed in concert
to ban any mining application, without any exceptions, such as, for example,
the fact that 'cpu mining' is one of the 4 primary functions this
application performs, and it is not deceptive in this respect.

With the shortage of CPUs, we anticipate that some people will want to
install the packages anyway, these non-`appstore` build tag versions will be
made available in binary releases alongside the rest of the versions, including
`headless` for fast deployment.

## Design goal of this project:

It is the belief of the ParallelCoin team that a blockchain network that has 
a 'centre' is defeating the whole purpose of the technology, and as such, 
the elimination of the utility of special purpose mining hardware is a key 
goal of this upgrade. 

The proof of work is based on long division, which 
cannot be accelerated any faster with custom silicon of a lower pitch than a 
modern CPU or GPU. There is currently no implementation of the hash function,
tentatively called DivHash, for GPU, however, it must be pointed out that 
being that GPUs (and mobile CPUs) have 32 bit wide dividers and thus have 
half the performance per clock cycle. Long division circuits are the largest 
single processing unit on CPU chips, and the procedure for calculating it is 
little different than how it is done manually by hand - requiring times 
tables, and proceeding one digit at a time from most to the least significant 
digit. 

Thus, the state-of-the-art technology for performing this hash 
function is already the CPU, and implicitly, could not possibly be made 
faster without the most advanced chip making hardware that is primarily only 
available for making CPUs due to economics.

The hash function jumbles and concatenates the raw block bytes, splits into 
two, squares each half, and then multiplies the halves, and finally divides 
by a jumbled version of the starting bytes. This is then repeated 3 times 
until a number is produced of the order of 20-40kb in size. It is possible 
to make this work target even harder, but beyond 3 cycles it starts to 
demand so much memory transfer that runs lower than one hash per second 
and 5 cycles produces nearly half a megabyte of data at the end. 

The result of the calculation is then hashed using the fast Blake3 hash 
function, to yield the result. Since there is no realistic way in which with 
such large numbers it could ever be shortcut, to produce the output value 
that goes through the hash function, this proof of work should permanently 
remain impossible to accelerate.

This is very important to containing the tendency towards centralisation, as 
on the whole, CPUs can perform the calculations at a rate that has not 
accelerated much faster than the raising of clock speeds, and thus a current 
model CPU, compared to one 4 years old, is not so dramatically faster, 
either. 

Most of the performance improvements come from larger amounts of on-chip 
cache memory, which is also the most expensive, and fastest memory there is. 
Without any easy way to acquire large numbers of processors to perform the 
work, without heavily competing with a currently (April 2021) very 
overstretched chip manufacturing industry.

It should ensure that no single miner, or even, as is the case with BTC and 
ETH, 20-30 miners, can dominate the mining, and thus threaten the security 
and stability of the network through either monopolistic practices or the 
incursion of government agencies into the business, in any one country, as 
the bulk of usable processors are distributed pretty evenly across the world in 
proportion with the relative size of the national economies.

The miner's primary setup, contrary to standard designs, uses a multicast 
gossip protocol to deliver new block templates to miner workers, and the 
sending of solutions back to the nodes, based on a pre-shared key for basic 
symmetric encryption security.
