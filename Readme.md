# vcardenc
by Cathal Garvey, copyright 2016, released under GNU AGPLv3 or later.

### About
This is my fourth or so try at an actual vCard parser/generator library in Go.
I attribute this largely to the terribleness of the vCard/iCal format, which
should probably be burned in entirety and replaced with either hCard or a sane
JSON specification.

So, if you're planning to use this to somehow broaden the reach of vCard or
further entrench its undeserved place in contacts and calendar interchange, don't.

Instead, use this to maintain, badly, legacy interchange and weep a little every
time you do.

I make no guarantees that this works correctly, in fact at present it's not even
passing the haphazard test cases I've put together.

### Goals
1. To parse vCard 4.0 to enable a rescue-path for contact data stuck in vCard.
2. To emit vCard 4.0 to provide a way to get sane contact data into systems that
   only accept vCard input, like many contact management things, or to create
   QR codes for business cards, etcetera.

### Not Goals
1. To properly implement the terribleness of the many vCard RFCs
2. To make this useful for anything but a crude rescue path
3. To deal with edge cases. For those, manually transcribe the data or prune
   stuff that breaks this parser.

### Status
1. Metadata parsing is broken on quoted data, as it tries to break quoted
   values on commas anyway. Needs an overhaul/refactor
2. Code is spaghettiish in many places and needs a refactor and more functionalisation.
3. API probably looks hideous on Godoc right now.
4. Linewise parsing is mostly complete but assumes lines have been unwrapped.
5. Card-wise parsing not even begun, yet.
